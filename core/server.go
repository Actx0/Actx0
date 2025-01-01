// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package core

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/clivern/actx0/api"
	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/middleware"
	"github.com/clivern/actx0/module"
	"github.com/clivern/actx0/pkg/ai"
	"github.com/clivern/actx0/pkg/qdrant"
	"github.com/clivern/actx0/pkg/storage"
	"github.com/clivern/actx0/service/agent"
	"github.com/clivern/actx0/service/knowledge"
	"github.com/clivern/actx0/task"

	"github.com/go-chi/chi/v5"
	cmid "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Setup creates and configures the HTTP server
func SetupServer(Static embed.FS) http.Handler {
	r := chi.NewRouter()

	r.Use(cmid.Recoverer)
	r.Use(cmid.Timeout(
		time.Duration(viper.GetInt("app.timeout")) * time.Second,
	))
	r.Use(middleware.RequestId)
	r.Use(middleware.AppContext)
	r.Use(middleware.PrometheusMiddleware)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestSizeLimit(
		int64(10 * 1024 * 1024),
	))
	r.Use(middleware.SessionAuth())
	r.Get("/favicon.ico", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	r.Group(func(r chi.Router) {
		r.Get("/api/v1/public/_health", api.HealthAction)
		r.Get("/api/v1/public/_ready", api.ReadyAction)
		r.Post("/api/v1/public/action/setup", api.SetupAction)
		r.Get("/api/v1/public/action/setup/status", api.SetupStatusAction)
		r.Post("/api/v1/public/action/login", api.LoginAction)
		r.Post("/api/v1/public/action/logout", api.LogoutAction)
		r.Post("/api/v1/public/action/forgot-password", api.ForgotPasswordAction)
		r.Post("/api/v1/public/action/reset-password", api.ResetPasswordAction)
		r.Post("/api/v1/public/action/register", api.RegisterAction)
		r.Post("/api/v1/public/action/stripe/webhook", api.StripeWebhookAction)
	})
	r.Get("/api/v1/me", api.GetMeAction)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireUser())
		r.Use(middleware.RequireRole(db.UserRoleAdmin, db.UserRoleRegular))
		r.Get("/api/v1/action/profile", api.GetProfileAction)
		r.Put("/api/v1/action/profile", api.UpdateProfileAction)
		r.Get("/api/v1/action/invites", api.ListUserInvitesAction)
		r.Get("/api/v1/action/invite-by-token/{token}", api.GetAuthenticatedUserInviteByTokenAction)
		r.Post("/api/v1/action/accept-invite/{token}", api.AcceptUserInviteByTokenAction)
		r.Post("/api/v1/action/reject-invite/{token}", api.RejectUserInviteByTokenAction)
	})
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireUser())
		r.Use(middleware.RequireRole(db.UserRoleAdmin))
		r.Put("/api/v1/action/settings", api.UpdateSettingsAction)
		r.Get("/api/v1/action/settings", api.GetSettingsAction)
	})
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireUser())
		r.Use(middleware.RequireRole(db.UserRoleAdmin, db.UserRoleRegular))
		r.Post("/api/v1/apiKeys", api.CreateUserAPIKeyAction)
		r.Get("/api/v1/apiKeys", api.ListUserAPIKeysAction)
		r.Get("/api/v1/apiKeys/{apiKeyId}", api.GetUserAPIKeyAction)
		r.Delete("/api/v1/apiKeys/{apiKeyId}", api.DeleteUserAPIKeyAction)
	})
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireUser())
		r.Use(middleware.RequireRole(db.UserRoleAdmin, db.UserRoleRegular))
		r.Post("/api/v1/workspaces", api.CreateWorkspaceAction)
		r.Get("/api/v1/workspaces", api.ListWorkspacesAction)
	})

	r.Route("/api/v1/workspaces/{workspaceId}", func(r chi.Router) {
		r.Use(middleware.RequireWorkspace())

		r.With(middleware.UserWorkspacePerm(module.CanGetWorkspace)...).Get("/", api.GetWorkspaceAction)
		r.With(middleware.UserWorkspacePerm(module.CanUpdateWorkspace)...).Put("/", api.UpdateWorkspaceAction)
		r.With(middleware.UserWorkspacePerm(module.CanDeleteWorkspace)...).Delete("/", api.DeleteWorkspaceAction)

		r.Route("/invites", func(r chi.Router) {
			r.With(middleware.UserWorkspacePerm(module.CanInviteMember)...).Post("/", api.CreateInviteAction)
			r.With(middleware.UserWorkspacePerm(module.CanListWorkspaceInvites)...).Get("/", api.ListInvitesAction)
			r.With(middleware.UserWorkspacePerm(module.CanGetWorkspaceInvite)...).Get("/{inviteId}", api.GetInviteAction)
			r.With(middleware.UserWorkspacePerm(module.CanRemoveInvite)...).Delete("/{inviteId}", api.DeleteInviteAction)
		})

		r.Route("/members", func(r chi.Router) {
			r.With(middleware.UserWorkspacePerm(module.CanListWorkspaceMembers)...).Get("/", api.ListWorkspaceMembersAction)
			r.With(middleware.UserWorkspacePerm(module.CanUpdateMemberRole)...).Put("/{memberUserId}", api.UpdateWorkspaceMemberRoleAction)
			r.With(middleware.UserWorkspacePerm(module.CanRemoveMember)...).Delete("/{memberUserId}", api.DeleteWorkspaceMemberAction)
		})

		r.Route("/keys", func(r chi.Router) {
			r.With(middleware.UserWorkspacePerm(module.CanCreateWorkspaceAccessKey)...).Post("/", api.CreateWorkspaceAccessKeyAction)
			r.With(middleware.UserWorkspacePerm(module.CanListWorkspaceAccessKeys)...).Get("/", api.ListWorkspaceAccessKeysAction)
			r.With(middleware.UserWorkspacePerm(module.CanGetWorkspaceAccessKey)...).Get("/{keyId}", api.GetWorkspaceAccessKeyAction)
			r.With(middleware.UserWorkspacePerm(module.CanDeleteWorkspaceAccessKey)...).Delete("/{keyId}", api.DeleteWorkspaceAccessKeyAction)
		})

		r.With(middleware.RequireWorkspacePermission(module.CanGetWorkspaceBilling)).Get("/billing", api.GetBillingStatusAction)
		r.With(middleware.UserWorkspacePerm(module.CanUpdateWorkspaceBilling)...).Post("/billing/checkout", api.CreateBillingCheckoutAction)
		r.With(middleware.UserWorkspacePerm(module.CanUpdateWorkspaceBilling)...).Post("/billing/portal", api.CreateBillingPortalAction)

		r.With(middleware.RequireWorkspacePermission(module.CanGetWorkspace)).Get("/stats", api.GetWorkspaceStatsAction)

		r.Route("/prompts", func(r chi.Router) {
			r.With(middleware.RequireWorkspacePermission(module.CanListPrompts)).Get("/", api.ListPromptsAction)
			r.With(middleware.UserWorkspacePerm(module.CanCreatePrompt)...).Post("/", api.CreatePromptAction)
			r.With(middleware.RequireWorkspacePermission(module.CanGetPrompt)).Get("/{promptId}", api.GetPromptAction)
			r.With(middleware.UserWorkspacePerm(module.CanDeletePrompt)...).Delete("/{promptId}", api.DeletePromptAction)

			r.Route("/{promptId}/versions", func(r chi.Router) {
				r.With(middleware.RequireWorkspacePermission(module.CanListPrompts)).Get("/", api.ListPromptVersionsAction)
				r.With(middleware.UserWorkspacePerm(module.CanCreatePrompt)...).Post("/", api.CreatePromptVersionAction)
				r.With(middleware.RequireWorkspacePermission(module.CanGetPrompt)).Get("/{promptVersionId}", api.GetPromptVersionAction)
				r.With(middleware.UserWorkspacePerm(module.CanUpdatePrompt)...).Put("/{promptVersionId}", api.UpdatePromptVersionAction)
				r.With(middleware.UserWorkspacePerm(module.CanDeletePrompt)...).Delete("/{promptVersionId}", api.DeletePromptVersionAction)
			})
		})

		r.With(middleware.RequireWorkspacePermission(module.CanGetPrompt)).Get("/promptsByName/{promptName}", api.GetPromptByNameAction)

		r.Route("/documents", func(r chi.Router) {
			r.With(middleware.RequireWorkspacePermission(module.CanListWorkspaceDocuments)).Get("/", api.ListDocumentsAction)
			r.With(middleware.RequireWorkspacePermission(module.CanQueryWorkspaceDocuments)).Post("/search", api.SearchDocumentsAction)
			r.With(middleware.UserWorkspacePerm(module.CanCreateWorkspaceDocument)...).Post("/", api.UploadDocumentAction)
			r.With(middleware.UserWorkspacePerm(module.CanDeleteWorkspaceDocument)...).Delete("/{documentId}", api.DeleteDocumentAction)
		})

		r.Route("/agents", func(r chi.Router) {
			r.With(middleware.RequireWorkspacePermission(module.CanListAgents)).Get("/", api.ListAgentsAction)
			r.With(middleware.UserWorkspacePerm(module.CanCreateAgent)...).Post("/", api.CreateAgentAction)
			r.With(middleware.RequireWorkspacePermission(module.CanGetAgent)).Get("/{agentId}", api.GetAgentAction)
			r.With(middleware.UserWorkspacePerm(module.CanUpdateAgent)...).Put("/{agentId}", api.UpdateAgentAction)
			r.With(middleware.UserWorkspacePerm(module.CanDeleteAgent)...).Delete("/{agentId}", api.DeleteAgentAction)

			r.Route("/{agentId}/sessions", func(r chi.Router) {
				r.With(middleware.RequireWorkspacePermission(module.CanListAgentSessions)).Get("/", api.ListAgentSessionsAction)
				r.With(middleware.RequireWorkspacePermission(module.CanGetAgentSession)).Get("/by-labels", api.GetAgentSessionByLabelsAction)
				r.With(middleware.RequireWorkspacePermission(module.CanGetAgentSession)).Get("/{sessionId}", api.GetAgentSessionAction)
				r.With(middleware.RequireWorkspacePermission(module.CanCreateAgentSession)).Post("/", api.CreateAgentSessionAction)
				r.With(middleware.RequireWorkspacePermission(module.CanUpdateAgentSession)).Put("/by-labels", api.UpdateAgentSessionByLabelsAction)
				r.With(middleware.RequireWorkspacePermission(module.CanDeleteAgentSession)).Delete("/by-labels", api.DeleteAgentSessionByLabelsAction)

				r.Route("/{sessionId}/messages", func(r chi.Router) {
					r.With(middleware.RequireWorkspacePermission(module.CanListSessionMessages)).Get("/", api.ListSessionMessagesAction)
					r.With(middleware.UserWorkspacePerm(module.CanCreateSessionMessage)...).Post("/", api.CreateSessionMessageAction)
					r.With(middleware.RequireWorkspacePermission(module.CanGetSessionMessage)).Get("/{messageId}", api.GetSessionMessageAction)
					r.With(middleware.UserWorkspacePerm(module.CanUpdateSessionMessage)...).Put("/{messageId}", api.UpdateSessionMessageAction)
					r.With(middleware.UserWorkspacePerm(module.CanDeleteSessionMessage)...).Delete("/{messageId}", api.DeleteSessionMessageAction)
				})

				r.Route("/{sessionId}/memories", func(r chi.Router) {
					r.With(middleware.RequireWorkspacePermission(module.CanListSessionMemories)).Get("/", api.ListSessionMemoriesAction)
					r.With(middleware.UserWorkspacePerm(module.CanCreateSessionMemory)...).Post("/", api.CreateSessionMemoryAction)
					r.With(middleware.RequireWorkspacePermission(module.CanGetSessionMemory)).Get("/{memoryId}", api.GetSessionMemoryAction)
					r.With(middleware.UserWorkspacePerm(module.CanUpdateSessionMemory)...).Put("/{memoryId}", api.UpdateSessionMemoryAction)
					r.With(middleware.UserWorkspacePerm(module.CanDeleteSessionMemory)...).Delete("/{memoryId}", api.DeleteSessionMemoryAction)
				})
			})
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireUser())
		r.Use(middleware.RequireRole(db.UserRoleAdmin, db.UserRoleRegular))
		r.Post("/api/v1/tool/getCityInformation", api.GetCityInformationAction)
		r.Post("/api/v1/tool/getWeatherByCity", api.GetWeatherByCityAction)
		r.Post("/api/v1/tool/getWeatherByCoordinates", api.GetWeatherByCoordinatesAction)
	})

	r.With(middleware.BasicAuth(
		viper.GetString("app.metrics.username"),
		viper.GetString("app.metrics.secret"),
	)).Get(
		"/api/v1/public/_metrics", promhttp.Handler().ServeHTTP,
	)

	dist, err := fs.Sub(Static, "web/dist")
	if err != nil {
		panic(fmt.Sprintf(
			"Error while accessing dist files: %s",
			err.Error(),
		))
	}

	r.Handle("/assets/*", http.StripPrefix("/", http.FileServer(http.FS(dist))))
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		indexFile, err := dist.Open("index.html")
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		defer indexFile.Close()

		stat, err := indexFile.Stat()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeContent(w, r, "index.html", stat.ModTime(), indexFile.(io.ReadSeeker))
	})

	return r
}

// InitDatabase initializes the database connection from configuration
func InitDatabase() error {
	rwConfig, roConfigs := DatabaseConfigs()

	return db.InitDB(rwConfig, roConfigs...)
}

// Run starts the HTTP server with graceful shutdown support
func RunServer(handler http.Handler) error {
	err := InitDatabase()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	module.RegisterEventListeners()

	asyr, err := module.Start(db.NewAsyncTaskRepository(db.GetDB(true)))
	if err != nil {
		return fmt.Errorf("failed to start async worker pool: %w", err)
	}

	store, err := storage.New()
	if err != nil {
		return fmt.Errorf("failed to initialize document storage: %w", err)
	}

	vdb, err := qdrant.New()
	if err != nil {
		return fmt.Errorf("failed to initialize qdrant: %w", err)
	}

	defer func() {
		err := vdb.Close()
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error closing qdrant client")
		}
	}()

	ksvc := knowledge.New(knowledge.Dependencies{
		Documents: db.NewWorkspaceDocumentRepository(
			db.GetDB(true),
		),
		Embed:         ai.NewEmbedClient(),
		Vectors:       vdb,
		Store:         store,
		Usage:         db.NewUsageRepository(db.GetDB(false)),
		Subscriptions: db.NewSubscriptionRepository(db.GetDB(false)),
	})

	asvc := agent.New(agent.Dependencies{
		Embed:    ai.NewEmbedClient(),
		Vectors:  vdb,
		LiteLLM:  ai.NewLiteClient(),
		LargeLLM: ai.NewLargeClient(),
		Agents:   db.NewAgentRepository(db.GetDB(true)),
		Sessions: db.NewAgentSessionRepository(db.GetDB(true)),
		Messages: db.NewSessionMessageRepository(db.GetDB(true)),
		Memories: db.NewSessionMemoryRepository(db.GetDB(true)),
	})

	task.Register(asyr, task.Dependencies{
		Knowledge: ksvc,
		Agent:     asvc,
	})

	defer func() {
		err := db.CloseDB()
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error closing database connection")
		}
	}()

	defer func() {
		err := asyr.Stop(30 * time.Second)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error stopping async worker pool")
		}
	}()

	timeout := time.Duration(viper.GetInt("app.timeout")) * time.Second

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", strconv.Itoa(viper.GetInt("app.port"))),
		Handler:           handler,
		MaxHeaderBytes:    230 * 1024,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadTimeout:       timeout + 5*time.Second,
		WriteTimeout:      timeout + 5*time.Second,
	}

	serr := make(chan error, 1)

	go func() {
		log.Info().
			Int("port", viper.GetInt("app.port")).
			Msg("Starting HTTP server")

		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			serr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serr:
		return fmt.Errorf("server error: %w", err)
	case sig := <-quit:
		log.Info().
			Str("signal", sig.String()).
			Msg("Received shutdown signal")

		shutdownTimeout := 30 * time.Second

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		log.Info().
			Dur("timeout", shutdownTimeout).
			Msg("Gracefully shutting down server")

		err := srv.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("server forced to shutdown: %w", err)
		}

		log.Info().Msg("Server shutdown complete")
	}

	return nil
}
