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
	EventUserLoggedIn       = "user.logged_in"
	EventUserLoggedOut      = "user.logged_out"
	EventUserForgotPassword = "user.forgot_password"
	EventUserResetPassword  = "user.reset_password"
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

// UserForgotPasswordData is emitted after a password reset token is created.
type UserForgotPasswordData struct {
	UserID db.Id
}

// UserResetPasswordData is emitted after a password is reset.
type UserResetPasswordData struct {
	UserID db.Id
}

var (
	UserLoggedIn       = event.New[UserLoggedInData](EventUserLoggedIn)
	UserLoggedOut      = event.New[UserLoggedOutData](EventUserLoggedOut)
	UserForgotPassword = event.New[UserForgotPasswordData](EventUserForgotPassword)
	UserResetPassword  = event.New[UserResetPasswordData](EventUserResetPassword)
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

	UserForgotPassword.On(func(_ context.Context, e UserForgotPasswordData) {
		log.Info().
			Str("userId", e.UserID.String()).
			Msg("Password reset token created")
	})

	UserResetPassword.On(func(_ context.Context, e UserResetPasswordData) {
		log.Info().
			Str("userId", e.UserID.String()).
			Msg("Password reset completed")
	})
}
