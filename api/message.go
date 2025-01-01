// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"errors"
	"net/http"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/locale"
	"github.com/clivern/actx0/module"
	"github.com/clivern/actx0/pkg/util"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// CreateSessionMessageAction creates a message in an agent session.
func CreateSessionMessageAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	agentIdStr := chi.URLParam(r, "agentId")
	if lo.IsEmpty(agentIdStr) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "invalid_agent_id"),
		})
		return
	}

	sessionIdStr := chi.URLParam(r, "sessionId")
	if lo.IsEmpty(sessionIdStr) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "invalid_session_reference"),
		})
		return
	}

	var req module.CreateSessionMessageRequest
	if err := util.DecodeAndValidate(r, &req); err != nil {
		util.WriteValidationError(w, err)
		return
	}

	mm := module.NewMessage(
		db.NewSessionMessageRepository(db.GetDB()),
		db.NewAgentSessionRepository(db.GetDB()),
	)
	message, err := mm.CreateSessionMessage(
		workspaceId,
		db.Id(agentIdStr),
		db.Id(sessionIdStr),
		&req,
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAgentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "agent_not_found"),
			})
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		default:
			log.Error().Err(err).Str("sessionId", sessionIdStr).Msg("Failed to create session message")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"errorMessage": locale.TR(r, "failed_create_session_message"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusCreated, message)
}

// BatchCreateSessionMessagesAction creates multiple messages in an agent session.
func BatchCreateSessionMessagesAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	agentIdStr := chi.URLParam(r, "agentId")
	if lo.IsEmpty(agentIdStr) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "invalid_agent_id"),
		})
		return
	}

	sessionIdStr := chi.URLParam(r, "sessionId")
	if lo.IsEmpty(sessionIdStr) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "invalid_session_reference"),
		})
		return
	}

	var req module.BatchCreateSessionMessagesRequest
	if err := util.DecodeAndValidate(r, &req); err != nil {
		util.WriteValidationError(w, err)
		return
	}

	mm := module.NewMessage(
		db.NewSessionMessageRepository(db.GetDB()),
		db.NewAgentSessionRepository(db.GetDB()),
	)
	messages, err := mm.BatchCreateSessionMessages(
		workspaceId,
		db.Id(agentIdStr),
		db.Id(sessionIdStr),
		&req,
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAgentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "agent_not_found"),
			})
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		default:
			log.Error().Err(err).Str("sessionId", sessionIdStr).Msg("Failed to batch create session messages")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"errorMessage": locale.TR(r, "failed_create_session_message"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"messages": messages,
	})
}

// ListSessionMessagesAction lists messages for an agent session.
func ListSessionMessagesAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	agentIdStr := chi.URLParam(r, "agentId")
	if lo.IsEmpty(agentIdStr) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "invalid_agent_id"),
		})
		return
	}

	sessionIdStr := chi.URLParam(r, "sessionId")
	if lo.IsEmpty(sessionIdStr) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "invalid_session_reference"),
		})
		return
	}

	limit, offset := util.ParsePagination(r)

	mm := module.NewMessage(
		db.NewSessionMessageRepository(db.GetDB()),
		db.NewAgentSessionRepository(db.GetDB()),
	)
	result, err := mm.ListSessionMessages(
		workspaceId,
		db.Id(agentIdStr),
		db.Id(sessionIdStr),
		limit,
		offset,
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAgentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "agent_not_found"),
			})
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		default:
			log.Error().Err(err).Str("sessionId", sessionIdStr).Msg("Failed to list session messages")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"errorMessage": locale.TR(r, "failed_list_session_messages"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"messages": result.Messages,
		"_meta": map[string]interface{}{
			"limit":  limit,
			"offset": offset,
			"total":  result.Total,
		},
	})
}

// GetSessionMessageAction returns one message by id.
func GetSessionMessageAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	agentIdStr := chi.URLParam(r, "agentId")
	if lo.IsEmpty(agentIdStr) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "invalid_agent_id"),
		})
		return
	}

	sessionIdStr := chi.URLParam(r, "sessionId")
	if lo.IsEmpty(sessionIdStr) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "invalid_session_reference"),
		})
		return
	}

	messageIdStr := chi.URLParam(r, "messageId")
	if lo.IsEmpty(messageIdStr) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "invalid_message_id"),
		})
		return
	}

	mm := module.NewMessage(
		db.NewSessionMessageRepository(db.GetDB()),
		db.NewAgentSessionRepository(db.GetDB()),
	)
	message, err := mm.GetSessionMessage(
		workspaceId,
		db.Id(agentIdStr),
		db.Id(sessionIdStr),
		db.Id(messageIdStr),
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAgentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "agent_not_found"),
			})
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		case errors.Is(err, module.ErrSessionMessageNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "session_message_not_found"),
			})
		default:
			log.Error().Err(err).Str("messageId", messageIdStr).Msg("Failed to get session message")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"errorMessage": locale.TR(r, "failed_get_session_message"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, message)
}

// UpdateSessionMessageAction updates a message by id.
func UpdateSessionMessageAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	agentIdStr := chi.URLParam(r, "agentId")
	if lo.IsEmpty(agentIdStr) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "invalid_agent_id"),
		})
		return
	}

	sessionIdStr := chi.URLParam(r, "sessionId")
	if lo.IsEmpty(sessionIdStr) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "invalid_session_reference"),
		})
		return
	}

	messageIdStr := chi.URLParam(r, "messageId")
	if lo.IsEmpty(messageIdStr) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "invalid_message_id"),
		})
		return
	}

	var req module.UpdateSessionMessageRequest
	if err := util.DecodeAndValidate(r, &req); err != nil {
		util.WriteValidationError(w, err)
		return
	}

	mm := module.NewMessage(
		db.NewSessionMessageRepository(db.GetDB()),
		db.NewAgentSessionRepository(db.GetDB()),
	)
	message, err := mm.UpdateSessionMessage(
		workspaceId,
		db.Id(agentIdStr),
		db.Id(sessionIdStr),
		db.Id(messageIdStr),
		&req,
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAgentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "agent_not_found"),
			})
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		case errors.Is(err, module.ErrSessionMessageNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "session_message_not_found"),
			})
		default:
			log.Error().Err(err).Str("messageId", messageIdStr).Msg("Failed to update session message")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"errorMessage": locale.TR(r, "failed_update_session_message"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, message)
}

// DeleteSessionMessageAction deletes a message by id.
func DeleteSessionMessageAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	agentIdStr := chi.URLParam(r, "agentId")
	if lo.IsEmpty(agentIdStr) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "invalid_agent_id"),
		})
		return
	}

	sessionIdStr := chi.URLParam(r, "sessionId")
	if lo.IsEmpty(sessionIdStr) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "invalid_session_reference"),
		})
		return
	}

	messageIdStr := chi.URLParam(r, "messageId")
	if lo.IsEmpty(messageIdStr) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "invalid_message_id"),
		})
		return
	}

	mm := module.NewMessage(
		db.NewSessionMessageRepository(db.GetDB()),
		db.NewAgentSessionRepository(db.GetDB()),
	)
	err := mm.DeleteSessionMessage(
		workspaceId,
		db.Id(agentIdStr),
		db.Id(sessionIdStr),
		db.Id(messageIdStr),
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAgentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "agent_not_found"),
			})
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		case errors.Is(err, module.ErrSessionMessageNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "session_message_not_found"),
			})
		default:
			log.Error().Err(err).Str("messageId", messageIdStr).Msg("Failed to delete session message")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"errorMessage": locale.TR(r, "failed_delete_session_message"),
			})
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// BatchDeleteSessionMessagesAction deletes multiple messages in an agent session.
func BatchDeleteSessionMessagesAction(w http.ResponseWriter, r *http.Request) {
	wid := chi.URLParam(r, "workspaceId")
	if lo.IsEmpty(wid) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "invalid_workspace_id"),
		})
		return
	}
	workspaceId := db.Id(wid)

	agentIdStr := chi.URLParam(r, "agentId")
	if lo.IsEmpty(agentIdStr) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "invalid_agent_id"),
		})
		return
	}

	sessionIdStr := chi.URLParam(r, "sessionId")
	if lo.IsEmpty(sessionIdStr) {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "invalid_session_reference"),
		})
		return
	}

	var req module.BatchDeleteSessionMessagesRequest
	if err := util.DecodeAndValidate(r, &req); err != nil {
		util.WriteValidationError(w, err)
		return
	}

	mm := module.NewMessage(
		db.NewSessionMessageRepository(db.GetDB()),
		db.NewAgentSessionRepository(db.GetDB()),
	)
	err := mm.BatchDeleteSessionMessages(
		workspaceId,
		db.Id(agentIdStr),
		db.Id(sessionIdStr),
		&req,
	)

	if err != nil {
		switch {
		case errors.Is(err, module.ErrWorkspaceNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "workspace_not_found"),
			})
		case errors.Is(err, module.ErrAgentNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "agent_not_found"),
			})
		case errors.Is(err, module.ErrAgentSessionNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "agent_session_not_found"),
			})
		case errors.Is(err, module.ErrSessionMessageNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": locale.TR(r, "session_message_not_found"),
			})
		default:
			log.Error().Err(err).Str("sessionId", sessionIdStr).Msg("Failed to batch delete session messages")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"errorMessage": locale.TR(r, "failed_delete_session_message"),
			})
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
