// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"context"
	"errors"
	"time"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/pkg/util"
)

// SessionManager manages sessions (create, revoke, etc.).
type SessionManager struct {
	SessionRepo db.SessionRepository
	UserRepo    db.UserRepository
}

// NewSessionManager creates a new session manager.
func NewSessionManager(sesr db.SessionRepository, ur db.UserRepository) *SessionManager {
	return &SessionManager{
		SessionRepo: sesr,
		UserRepo:    ur,
	}
}

// CreateSession creates a new session for a user.
func (s *SessionManager) CreateSession(ctx context.Context, userId db.Id, duration time.Duration) (*db.Session, error) {
	user, err := s.UserRepo.GetById(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	token, err := util.GenerateSecureToken(32)
	if err != nil {
		return nil, err
	}

	session := &db.Session{
		Token:     token,
		UserId:    userId,
		ExpiresAt: time.Now().UTC().Add(duration),
	}
	if appCtx, ok := AppCtxFrom(ctx); ok && appCtx != nil {
		if appCtx.IPAddress != "" {
			session.IPAddress = &appCtx.IPAddress
		}
		if appCtx.UserAgent != "" {
			session.UserAgent = &appCtx.UserAgent
		}
	}

	err = s.SessionRepo.Create(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// ValidateSession validates a session token and returns the associated user.
func (s *SessionManager) ValidateSession(token string) (*db.User, *db.Session, error) {
	session, err := s.SessionRepo.GetByToken(token)
	if err != nil {
		return nil, nil, err
	}
	if session == nil {
		return nil, nil, errors.New("session not found")
	}
	if session.ExpiresAt.Before(time.Now().UTC()) {
		s.SessionRepo.Delete(session.Id)
		return nil, nil, errors.New("session expired")
	}

	user, err := s.UserRepo.GetById(session.UserId)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, errors.New("user not found")
	}
	if !user.IsActive {
		return nil, nil, errors.New("user is not active")
	}

	return user, session, nil
}

// RevokeUserSessions revokes all sessions for a user.
func (s *SessionManager) RevokeUserSessions(userId db.Id) error {
	return s.SessionRepo.DeleteByUserId(userId)
}

// GetUserSessions retrieves all active sessions for a user.
func (s *SessionManager) GetUserSessions(userId db.Id) ([]*db.Session, error) {
	sessions, err := s.SessionRepo.GetByUserId(userId)
	if err != nil {
		return nil, err
	}

	var activeSessions []*db.Session
	now := time.Now().UTC()
	for _, session := range sessions {
		if session.ExpiresAt.After(now) {
			activeSessions = append(activeSessions, session)
		}
	}

	return activeSessions, nil
}

// CleanupExpiredSessions removes all expired sessions from the database.
func (s *SessionManager) CleanupExpiredSessions() (int64, error) {
	return s.SessionRepo.DeleteExpired()
}
