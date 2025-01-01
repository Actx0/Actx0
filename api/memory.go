// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"errors"
	"net/http"

	"github.com/clivern/actx0/conf"
	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/locale"
	"github.com/clivern/actx0/module"
	"github.com/clivern/actx0/pkg/ai"
	"github.com/clivern/actx0/pkg/qdrant"
	"github.com/clivern/actx0/pkg/util"
	"github.com/clivern/actx0/service/agent"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// CreateSessionMemoryAction creates a memory in an agent session.
func CreateSessionMemoryAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	agentId := chi.URLParam(r, "agentId")
	if lo.IsEmpty(agentId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_agent_id"),
		})
		return
	}

	sessionId := chi.URLParam(r, "sessionId")
	if lo.IsEmpty(sessionId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_session_reference"),
		})
		return
	}

	var req module.CreateSessionMemoryRequest
	if err := util.DecodeAndValidate(r, &req); err != nil {
		util.WriteValidationError(w, err)
		return
	}

	mm := module.NewMemory(
		db.NewSessionMemoryRepository(db.GetDB()),
		db.NewSessionMessageRepository(db.GetDB()),
		module.NewAgentSession(
			db.NewAgentSessionRepository(db.GetDB()),
			db.NewAgentRepository(db.GetDB()),
			db.NewWorkspaceRepository(db.GetDB()),
		),
	)
	memory, err := mm.CreateSessionMemory(
		workspaceId,
		db.Id(agentId),
		db.Id(sessionId),
		&req,
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAgentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_not_found"),
			})
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		default:
			log.Error().Err(err).Str("sessionId", sessionId).Msg("Failed to create session memory")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_create_session_memory"),
			})
		}
		return
	}

	defer record(r, &module.RecordAuditEvent{
		WorkspaceId:  workspaceId,
		Action:       ActionMemoryCreate,
		ResourceType: ResourceMemory,
		ResourceId:   memory.Id,
		Meta:         map[string]any{"sessionId": sessionId},
	})

	util.WriteJSON(w, http.StatusCreated, memory)
}

// BatchCreateSessionMemoriesAction creates multiple memories in an agent session.
func BatchCreateSessionMemoriesAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	agentId := chi.URLParam(r, "agentId")
	if lo.IsEmpty(agentId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_agent_id"),
		})
		return
	}

	sessionId := chi.URLParam(r, "sessionId")
	if lo.IsEmpty(sessionId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_session_reference"),
		})
		return
	}

	var req module.BatchCreateSessionMemoriesRequest
	if err := util.DecodeAndValidate(r, &req); err != nil {
		util.WriteValidationError(w, err)
		return
	}

	mm := module.NewMemory(
		db.NewSessionMemoryRepository(db.GetDB()),
		db.NewSessionMessageRepository(db.GetDB()),
		module.NewAgentSession(
			db.NewAgentSessionRepository(db.GetDB()),
			db.NewAgentRepository(db.GetDB()),
			db.NewWorkspaceRepository(db.GetDB()),
		),
	)
	memories, err := mm.BatchCreateSessionMemories(
		workspaceId,
		db.Id(agentId),
		db.Id(sessionId),
		&req,
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAgentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_not_found"),
			})
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		default:
			log.Error().Err(err).Str("sessionId", sessionId).Msg("Failed to batch create session memories")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_create_session_memory"),
			})
		}
		return
	}

	defer record(r, &module.RecordAuditEvent{
		WorkspaceId:  workspaceId,
		Action:       ActionMemoryBatchCreate,
		ResourceType: ResourceSession,
		ResourceId:   db.Id(sessionId),
		Meta:         map[string]any{"count": len(memories)},
	})

	util.WriteJSON(w, http.StatusCreated, map[string]any{
		"memories": memories,
	})
}

// ListSessionMemoriesAction lists memories for an agent session.
func ListSessionMemoriesAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	agentId := chi.URLParam(r, "agentId")
	if lo.IsEmpty(agentId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_agent_id"),
		})
		return
	}

	sessionId := chi.URLParam(r, "sessionId")
	if lo.IsEmpty(sessionId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_session_reference"),
		})
		return
	}

	limit, offset := util.ParsePagination(r)

	mm := module.NewMemory(
		db.NewSessionMemoryRepository(db.GetDB()),
		db.NewSessionMessageRepository(db.GetDB()),
		module.NewAgentSession(
			db.NewAgentSessionRepository(db.GetDB()),
			db.NewAgentRepository(db.GetDB()),
			db.NewWorkspaceRepository(db.GetDB()),
		),
	)
	result, err := mm.ListSessionMemories(
		workspaceId,
		db.Id(agentId),
		db.Id(sessionId),
		limit,
		offset,
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAgentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_not_found"),
			})
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		default:
			log.Error().Err(err).Str("sessionId", sessionId).Msg("Failed to list session memories")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_list_session_memories"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"memories": result.Memories,
		"_meta": map[string]any{
			"limit":  limit,
			"offset": offset,
			"total":  result.Total,
		},
	})
}

// GetSessionMemoryAction returns one memory by id.
func GetSessionMemoryAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	agentId := chi.URLParam(r, "agentId")
	if lo.IsEmpty(agentId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_agent_id"),
		})
		return
	}

	sessionId := chi.URLParam(r, "sessionId")
	if lo.IsEmpty(sessionId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_session_reference"),
		})
		return
	}

	memoryId := chi.URLParam(r, "memoryId")
	if lo.IsEmpty(memoryId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_memory_id"),
		})
		return
	}

	mm := module.NewMemory(
		db.NewSessionMemoryRepository(db.GetDB()),
		db.NewSessionMessageRepository(db.GetDB()),
		module.NewAgentSession(
			db.NewAgentSessionRepository(db.GetDB()),
			db.NewAgentRepository(db.GetDB()),
			db.NewWorkspaceRepository(db.GetDB()),
		),
	)
	memory, err := mm.GetSessionMemory(
		workspaceId,
		db.Id(agentId),
		db.Id(sessionId),
		db.Id(memoryId),
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAgentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_not_found"),
			})
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		case errors.Is(err, module.ErrSessionMemoryNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "session_memory_not_found"),
			})
		default:
			log.Error().Err(err).Str("memoryId", memoryId).Msg("Failed to get session memory")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_get_session_memory"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, memory)
}

// UpdateSessionMemoryAction updates a memory by id.
func UpdateSessionMemoryAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	agentId := chi.URLParam(r, "agentId")
	if lo.IsEmpty(agentId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_agent_id"),
		})
		return
	}

	sessionId := chi.URLParam(r, "sessionId")
	if lo.IsEmpty(sessionId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_session_reference"),
		})
		return
	}

	memoryId := chi.URLParam(r, "memoryId")
	if lo.IsEmpty(memoryId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_memory_id"),
		})
		return
	}

	var req module.UpdateSessionMemoryRequest
	if err := util.DecodeAndValidate(r, &req); err != nil {
		util.WriteValidationError(w, err)
		return
	}

	mm := module.NewMemory(
		db.NewSessionMemoryRepository(db.GetDB()),
		db.NewSessionMessageRepository(db.GetDB()),
		module.NewAgentSession(
			db.NewAgentSessionRepository(db.GetDB()),
			db.NewAgentRepository(db.GetDB()),
			db.NewWorkspaceRepository(db.GetDB()),
		),
	)
	memory, err := mm.UpdateSessionMemory(
		workspaceId,
		db.Id(agentId),
		db.Id(sessionId),
		db.Id(memoryId),
		&req,
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAgentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_not_found"),
			})
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		case errors.Is(err, module.ErrSessionMemoryNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "session_memory_not_found"),
			})
		default:
			log.Error().Err(err).Str("memoryId", memoryId).Msg("Failed to update session memory")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_update_session_memory"),
			})
		}
		return
	}

	defer record(r, &module.RecordAuditEvent{
		WorkspaceId:  workspaceId,
		Action:       ActionMemoryUpdate,
		ResourceType: ResourceMemory,
		ResourceId:   memory.Id,
		Meta:         map[string]any{"sessionId": sessionId},
	})

	util.WriteJSON(w, http.StatusOK, memory)
}

// DeleteSessionMemoryAction deletes a memory by id.
func DeleteSessionMemoryAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	agentId := chi.URLParam(r, "agentId")
	if lo.IsEmpty(agentId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_agent_id"),
		})
		return
	}

	sessionId := chi.URLParam(r, "sessionId")
	if lo.IsEmpty(sessionId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_session_reference"),
		})
		return
	}

	memoryId := chi.URLParam(r, "memoryId")
	if lo.IsEmpty(memoryId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_memory_id"),
		})
		return
	}

	mm := module.NewMemory(
		db.NewSessionMemoryRepository(db.GetDB()),
		db.NewSessionMessageRepository(db.GetDB()),
		module.NewAgentSession(
			db.NewAgentSessionRepository(db.GetDB()),
			db.NewAgentRepository(db.GetDB()),
			db.NewWorkspaceRepository(db.GetDB()),
		),
	)
	err := mm.DeleteSessionMemory(
		workspaceId,
		db.Id(agentId),
		db.Id(sessionId),
		db.Id(memoryId),
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAgentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_not_found"),
			})
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		case errors.Is(err, module.ErrSessionMemoryNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "session_memory_not_found"),
			})
		default:
			log.Error().Err(err).Str("memoryId", memoryId).Msg("Failed to delete session memory")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_delete_session_memory"),
			})
		}
		return
	}

	defer record(r, &module.RecordAuditEvent{
		WorkspaceId:  workspaceId,
		Action:       ActionMemoryDelete,
		ResourceType: ResourceMemory,
		ResourceId:   db.Id(memoryId),
		Meta:         map[string]any{"sessionId": sessionId},
	})

	w.WriteHeader(http.StatusNoContent)
}

// RegenerateSessionMemoriesAction deletes session memories and queues generation from user messages.
func RegenerateSessionMemoriesAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	agentId := chi.URLParam(r, "agentId")
	if lo.IsEmpty(agentId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_agent_id"),
		})
		return
	}

	sessionId := chi.URLParam(r, "sessionId")
	if lo.IsEmpty(sessionId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_session_reference"),
		})
		return
	}

	mm := module.NewMemory(
		db.NewSessionMemoryRepository(db.GetDB()),
		db.NewSessionMessageRepository(db.GetDB()),
		module.NewAgentSession(
			db.NewAgentSessionRepository(db.GetDB()),
			db.NewAgentRepository(db.GetDB()),
			db.NewWorkspaceRepository(db.GetDB()),
		),
	)
	result, err := mm.RegenerateSessionMemories(
		workspaceId,
		db.Id(agentId),
		db.Id(sessionId),
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAgentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_not_found"),
			})
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		default:
			log.Error().Err(err).Str("sessionId", sessionId).Msg("Failed to regenerate session memories")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_regenerate_session_memories"),
			})
		}
		return
	}

	defer record(r, &module.RecordAuditEvent{
		WorkspaceId:  workspaceId,
		Action:       ActionMemoryRegenerate,
		ResourceType: ResourceMemory,
		ResourceId:   db.Id(sessionId),
		Meta:         map[string]any{"sessionId": sessionId, "queued": result.Queued},
	})

	util.WriteJSON(w, http.StatusAccepted, result)
}

// BatchDeleteSessionMemoriesAction deletes multiple memories in an agent session.
func BatchDeleteSessionMemoriesAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	agentId := chi.URLParam(r, "agentId")
	if lo.IsEmpty(agentId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_agent_id"),
		})
		return
	}

	sessionId := chi.URLParam(r, "sessionId")
	if lo.IsEmpty(sessionId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_session_reference"),
		})
		return
	}

	var req module.BatchDeleteSessionMemoriesRequest
	if err := util.DecodeAndValidate(r, &req); err != nil {
		util.WriteValidationError(w, err)
		return
	}

	mm := module.NewMemory(
		db.NewSessionMemoryRepository(db.GetDB()),
		db.NewSessionMessageRepository(db.GetDB()),
		module.NewAgentSession(
			db.NewAgentSessionRepository(db.GetDB()),
			db.NewAgentRepository(db.GetDB()),
			db.NewWorkspaceRepository(db.GetDB()),
		),
	)
	err := mm.BatchDeleteSessionMemories(
		workspaceId,
		db.Id(agentId),
		db.Id(sessionId),
		&req,
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAgentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_not_found"),
			})
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		case errors.Is(err, module.ErrSessionMemoryNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "session_memory_not_found"),
			})
		default:
			log.Error().Err(err).Str("sessionId", sessionId).Msg("Failed to batch delete session memories")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_delete_session_memory"),
			})
		}
		return
	}

	defer record(r, &module.RecordAuditEvent{
		WorkspaceId:  workspaceId,
		Action:       ActionMemoryBatchDelete,
		ResourceType: ResourceSession,
		ResourceId:   db.Id(sessionId),
		Meta:         map[string]any{"count": len(req.IDs)},
	})

	w.WriteHeader(http.StatusNoContent)
}

// SearchSessionMemoriesAction searches session memories by semantic query.
func SearchSessionMemoriesAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	agentId := chi.URLParam(r, "agentId")
	if lo.IsEmpty(agentId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_agent_id"),
		})
		return
	}

	sessionId := chi.URLParam(r, "sessionId")
	if lo.IsEmpty(sessionId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_session_reference"),
		})
		return
	}

	var req module.SearchSessionMemoriesRequest
	if err := util.DecodeAndValidate(r, &req); err != nil {
		util.WriteValidationError(w, err)
		return
	}

	limit := lo.Ternary(req.Limit == 0, conf.DefaultSearchLimit, req.Limit)

	vdb, err := qdrant.New()
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize qdrant")
		util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"errorMessage": locale.TR(r, "failed_process_request"),
		})
		return
	}
	defer func() {
		if err := vdb.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close qdrant client")
		}
	}()

	asvc := agent.New(agent.Dependencies{
		Embed:    ai.NewEmbedClient(),
		Vectors:  vdb,
		Agents:   db.NewAgentRepository(db.GetDB()),
		Sessions: db.NewAgentSessionRepository(db.GetDB()),
		Messages: db.NewSessionMessageRepository(db.GetDB()),
		Memories: db.NewSessionMemoryRepository(db.GetDB()),
	})

	mm := module.NewMemory(
		db.NewSessionMemoryRepository(db.GetDB()),
		db.NewSessionMessageRepository(db.GetDB()),
		module.NewAgentSession(
			db.NewAgentSessionRepository(db.GetDB()),
			db.NewAgentRepository(db.GetDB()),
			db.NewWorkspaceRepository(db.GetDB()),
		),
	)
	mm.Agent = asvc
	mm.UsageRepository = db.NewUsageRepository(db.GetDB())
	mm.SubscriptionRepository = db.NewSubscriptionRepository(db.GetDB())

	result, err := mm.SearchSessionMemories(
		r.Context(),
		workspaceId,
		db.Id(agentId),
		db.Id(sessionId),
		req.Query,
		limit,
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAgentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_not_found"),
			})
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		default:
			log.Error().Err(err).Str("sessionId", sessionId).Msg("Failed to search session memories")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_search_session_memories"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, result)
}
