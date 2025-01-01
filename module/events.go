// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"context"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/pkg/event"

	"github.com/rs/zerolog/log"
)

const (
	EventUserLoggedIn  = "user.logged_in"
	EventUserLoggedOut = "user.logged_out"
)

// UserLoggedInData is emitted after a successful login.
type UserLoggedInData struct {
	UserID    db.Id
	SessionID db.Id
}

// UserLoggedOutData is emitted after a successful logout.
type UserLoggedOutData struct {
	UserID db.Id
}

var (
	UserLoggedIn  = event.New[UserLoggedInData](EventUserLoggedIn)
	UserLoggedOut = event.New[UserLoggedOutData](EventUserLoggedOut)
)

// RegisterEventListeners wires application event handlers.
func RegisterEventListeners() {
	UserLoggedIn.On(func(_ context.Context, e UserLoggedInData) {
		log.Info().
			Str("userId", e.UserID.String()).
			Str("sessionId", e.SessionID.String()).
			Msg("User logged in")
	})

	UserLoggedOut.On(func(_ context.Context, e UserLoggedOutData) {
		log.Info().
			Str("userId", e.UserID.String()).
			Msg("User logged out")
	})
}
