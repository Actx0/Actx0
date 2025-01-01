// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/pkg/util"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

var (
	ErrAuditEventNotFound    = errors.New("audit event not found")
	ErrFailedListAuditEvents = errors.New("failed list audit events")
	ErrFailedGetAuditEvent   = errors.New("failed get audit event")
)

// Audit actor types.
const (
	AuditActorUser      = "user"
	AuditActorAccessKey = "access_key"
)

// AuditEventResponse is an audit event shaped for API responses.
type AuditEventResponse struct {
	Id           db.Id           `json:"id"`
	WorkspaceId  db.Id           `json:"workspaceId"`
	ActorId      *db.Id          `json:"actorId,omitempty"`
	ActorName    *string         `json:"actorName,omitempty"`
	ActorType    *string         `json:"actorType,omitempty"`
	Action       string          `json:"action"`
	ResourceType *string         `json:"resourceType,omitempty"`
	ResourceId   *db.Id          `json:"resourceId,omitempty"`
	IPAddress    *string         `json:"ipAddress,omitempty"`
	UserAgent    *string         `json:"userAgent,omitempty"`
	Meta         json.RawMessage `json:"meta,omitempty"`
	CreatedAt    string          `json:"createdAt"`
}

// ListAuditEventsResponse is returned when listing audit events.
type ListAuditEventsResponse struct {
	Events []*AuditEventResponse
	Total  int64
}

// Audit is the module for workspace audit events.
type Audit struct {
	AuditRepository     db.AuditEventRepository
	WorkspaceRepository db.WorkspaceRepository
}

// NewAudit creates an audit module with the given repositories.
func NewAudit(audits db.AuditEventRepository, workspaces db.WorkspaceRepository) *Audit {
	return &Audit{
		AuditRepository:     audits,
		WorkspaceRepository: workspaces,
	}
}

// RecordAuditEvent describes an audit event to record.
type RecordAuditEvent struct {
	WorkspaceId  db.Id
	Action       string
	ResourceType string
	ResourceId   db.Id
	Meta         map[string]any
}

// Record persists an audit event, resolving the actor (user or access key)
func (a *Audit) Record(ctx context.Context, event *RecordAuditEvent) {
	row := &db.AuditEvent{
		WorkspaceId: event.WorkspaceId,
		Action:      event.Action,
	}

	if lo.IsNotEmpty(event.ResourceType) {
		row.ResourceType = &event.ResourceType
	}
	if lo.IsNotEmpty(event.ResourceId) {
		row.ResourceId = &event.ResourceId
	}

	if p, ok := PrincipalFrom(ctx); ok {
		switch {
		case p.IsUser():
			row.ActorId = &p.User.Id
			row.ActorName = &p.User.Name
			row.ActorType = new(AuditActorUser)
		case p.IsAccessKey():
			row.ActorId = &p.AccessKey.Id
			row.ActorName = &p.AccessKey.Name
			row.ActorType = new(AuditActorAccessKey)
		}
	}

	if appCtx, ok := AppCtxFrom(ctx); ok && appCtx != nil {
		if lo.IsNotEmpty(appCtx.IPAddress) {
			row.IPAddress = &appCtx.IPAddress
		}
		if lo.IsNotEmpty(appCtx.UserAgent) {
			row.UserAgent = &appCtx.UserAgent
		}
	}

	if len(event.Meta) > 0 {
		encoded, _ := json.Marshal(event.Meta)
		row.Meta = new(string(encoded))
	}

	err := a.AuditRepository.Create(row)
	if err != nil {
		log.Error().
			Err(err).
			Str("workspaceId", event.WorkspaceId.String()).
			Str("action", event.Action).
			Msg("Failed to record audit event")
	}
}

// GetAuditEvent returns one audit event by id.
func (a *Audit) GetAuditEvent(workspaceId, auditId db.Id) (*AuditEventResponse, error) {
	workspace, err := a.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	event, err := a.AuditRepository.GetById(auditId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetAuditEvent, err)
	}
	if event == nil || event.WorkspaceId != workspaceId {
		return nil, ErrAuditEventNotFound
	}

	return &AuditEventResponse{
		Id:           event.Id,
		WorkspaceId:  event.WorkspaceId,
		ActorId:      event.ActorId,
		ActorName:    event.ActorName,
		ActorType:    event.ActorType,
		Action:       event.Action,
		ResourceType: event.ResourceType,
		ResourceId:   event.ResourceId,
		IPAddress:    event.IPAddress,
		UserAgent:    event.UserAgent,
		Meta:         util.JSONRawFromString(event.Meta),
		CreatedAt:    event.CreatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// ListAuditEvents returns paginated audit events for a workspace.
func (a *Audit) ListAuditEvents(workspaceId db.Id, limit, offset int) (*ListAuditEventsResponse, error) {
	workspace, err := a.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	total, err := a.AuditRepository.CountByWorkspaceId(workspaceId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListAuditEvents, err)
	}

	events, err := a.AuditRepository.ListByWorkspaceId(workspaceId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListAuditEvents, err)
	}

	list := make([]*AuditEventResponse, 0, len(events))
	for _, item := range events {
		list = append(list, &AuditEventResponse{
			Id:           item.Id,
			WorkspaceId:  item.WorkspaceId,
			ActorId:      item.ActorId,
			ActorName:    item.ActorName,
			ActorType:    item.ActorType,
			Action:       item.Action,
			ResourceType: item.ResourceType,
			ResourceId:   item.ResourceId,
			IPAddress:    item.IPAddress,
			UserAgent:    item.UserAgent,
			Meta:         util.JSONRawFromString(item.Meta),
			CreatedAt:    item.CreatedAt.UTC().Format(time.RFC3339),
		})
	}

	return &ListAuditEventsResponse{Events: list, Total: total}, nil
}
