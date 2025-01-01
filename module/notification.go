// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"errors"
	"fmt"
	"time"

	"github.com/clivern/actx0/db"
)

var (
	ErrFailedListNotifications    = errors.New("failed list notifications")
	ErrFailedMarkNotificationRead = errors.New("failed mark notification read")
	ErrNotificationNotFound       = errors.New("notification not found")
)

// NotificationResponse is a notification shaped for API responses.
type NotificationResponse struct {
	Id          db.Id   `json:"id"`
	UserId      *db.Id  `json:"userId,omitempty"`
	WorkspaceId *db.Id  `json:"workspaceId,omitempty"`
	Type        string  `json:"type"`
	Title       string  `json:"title"`
	Body        *string `json:"body,omitempty"`
	ReadAt      *string `json:"readAt,omitempty"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

// ListNotificationsResponse is returned when listing notifications.
type ListNotificationsResponse struct {
	Notifications []*NotificationResponse
	Total         int64
}

// Notification is the module for user/workspace notifications.
type Notification struct {
	NotificationRepository db.NotificationRepository
	WorkspaceRepository    db.WorkspaceRepository
}

// NewNotification creates a notification module with the given repositories.
func NewNotification(notifications db.NotificationRepository, workspaces db.WorkspaceRepository) *Notification {
	return &Notification{
		NotificationRepository: notifications,
		WorkspaceRepository:    workspaces,
	}
}

// ListNotifications returns paginated inbox notifications for a user in a workspace.
func (n *Notification) ListNotifications(userId, workspaceId db.Id, limit, offset int) (*ListNotificationsResponse, error) {
	workspace, err := n.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	total, err := n.NotificationRepository.CountByUserIdAndWorkspaceId(userId, workspaceId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListNotifications, err)
	}

	rows, err := n.NotificationRepository.ListByUserIdAndWorkspaceId(userId, workspaceId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListNotifications, err)
	}

	list := make([]*NotificationResponse, 0, len(rows))
	for _, item := range rows {
		list = append(list, toNotificationResponse(item))
	}

	return &ListNotificationsResponse{Notifications: list, Total: total}, nil
}

// MarkNotificationRead marks a notification as read for the user's workspace inbox.
func (n *Notification) MarkNotificationRead(userId, workspaceId, notificationId db.Id) (*NotificationResponse, error) {
	workspace, err := n.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	item, err := n.NotificationRepository.GetById(notificationId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedMarkNotificationRead, err)
	}
	if item == nil {
		return nil, ErrNotificationNotFound
	}

	belongsToUser := item.UserId != nil && *item.UserId == userId
	belongsToWorkspace := item.WorkspaceId != nil && *item.WorkspaceId == workspaceId
	if !belongsToUser && !belongsToWorkspace {
		return nil, ErrNotificationNotFound
	}

	err = n.NotificationRepository.MarkRead(notificationId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedMarkNotificationRead, err)
	}

	item, err = n.NotificationRepository.GetById(notificationId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedMarkNotificationRead, err)
	}

	return toNotificationResponse(item), nil
}

func toNotificationResponse(item *db.Notification) *NotificationResponse {
	res := &NotificationResponse{
		Id:          item.Id,
		UserId:      item.UserId,
		WorkspaceId: item.WorkspaceId,
		Type:        item.Type,
		Title:       item.Title,
		Body:        item.Body,
		CreatedAt:   item.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   item.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if item.ReadAt != nil {
		readAt := item.ReadAt.UTC().Format(time.RFC3339)
		res.ReadAt = &readAt
	}
	return res
}
