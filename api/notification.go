// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"errors"
	"net/http"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/locale"
	"github.com/clivern/actx0/middleware"
	"github.com/clivern/actx0/module"
	"github.com/clivern/actx0/pkg/util"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// ListWorkspaceNotificationsAction lists inbox notifications for the current user in a workspace.
func ListWorkspaceNotificationsAction(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"errorMessage": locale.TR(r, "not_authenticated"),
		})
		return
	}

	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	limit, offset := util.ParsePagination(r)

	nm := module.NewNotification(
		db.NewNotificationRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	result, err := nm.ListNotifications(user.Id, workspaceId, limit, offset)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		default:
			log.Error().Err(err).Str("workspaceId", wid).Msg("Failed to list notifications")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_list_notifications"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"notifications": result.Notifications,
		"_meta": map[string]any{
			"limit":  limit,
			"offset": offset,
			"total":  result.Total,
		},
	})
}

// MarkWorkspaceNotificationReadAction marks a notification as read for the current user.
func MarkWorkspaceNotificationReadAction(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"errorMessage": locale.TR(r, "not_authenticated"),
		})
		return
	}

	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	notificationId := chi.URLParam(r, "notificationId")
	if lo.IsEmpty(notificationId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_notification_id"),
		})
		return
	}

	nm := module.NewNotification(
		db.NewNotificationRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	notification, err := nm.MarkNotificationRead(
		user.Id,
		db.Id(wid),
		db.Id(notificationId),
	)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrNotificationNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "notification_not_found"),
			})
		default:
			log.Error().Err(err).Str("notificationId", notificationId).Msg("Failed to mark notification as read")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_mark_notification_read"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, notification)
}

// MarkWorkspaceNotificationArchivedAction marks a notification as archived for the current user.
func MarkWorkspaceNotificationArchivedAction(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		util.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"errorMessage": locale.TR(r, "not_authenticated"),
		})
		return
	}

	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}

	notificationId := chi.URLParam(r, "notificationId")
	if lo.IsEmpty(notificationId) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"errorMessage": locale.TR(r, "invalid_notification_id"),
		})
		return
	}

	nm := module.NewNotification(
		db.NewNotificationRepository(db.GetDB()),
		db.NewWorkspaceRepository(db.GetDB()),
	)
	notification, err := nm.MarkNotificationArchived(
		user.Id,
		db.Id(wid),
		db.Id(notificationId),
	)
	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrNotificationNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "notification_not_found"),
			})
		default:
			log.Error().Err(err).Str("notificationId", notificationId).Msg("Failed to mark notification as archived")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_mark_notification_archived"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, notification)
}
