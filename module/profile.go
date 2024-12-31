// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/pkg/resend"
	"github.com/clivern/actx0/pkg/util"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

var (
	ErrCurrentPasswordRequired       = errors.New("current password required")
	ErrCurrentPasswordIncorrect      = errors.New("current password incorrect")
	ErrFailedGetProfile              = errors.New("failed get profile")
	ErrFailedUpdateProfile           = errors.New("failed update profile")
	ErrEmailAlreadyVerified          = errors.New("email already verified")
	ErrVerifyEmailResendLimitReached = errors.New("verify email resend daily limit reached")
	ErrFailedResendVerifyEmail       = errors.New("failed resend verify email")
)

const (
	VerifyEmailResendMetaKey    = "verify_email_resend"
	VerifyEmailResendDailyLimit = 5
)

// UpdateProfileRequest is the body for updating your profile.
type UpdateProfileRequest struct {
	Name            string `json:"name" validate:"omitempty,max=255" label:"Name"`
	Email           string `json:"email" validate:"omitempty,email,max=255" label:"Email"`
	CurrentPassword string `json:"currentPassword" validate:"omitempty,max=60" label:"Current password"`
	NewPassword     string `json:"newPassword" validate:"omitempty,strong_password,min=8,max=60" label:"New password"`
	Language        string `json:"language" validate:"required,oneof=en fr" label:"Language"`
	Theme           string `json:"theme" validate:"required,oneof=default blue slate emerald dark" label:"Theme"`
}

// Profile handles profile updates.
type Profile struct {
	UserRepository     db.UserRepository
	UserMetaRepository db.UserMetaRepository
}

// NewProfile creates a Profile module with the given repositories.
func NewProfile(users db.UserRepository, meta db.UserMetaRepository) *Profile {
	return &Profile{
		UserRepository:     users,
		UserMetaRepository: meta,
	}
}

// GetProfile retrieves a user's profile.
func (p *Profile) GetProfile(userId db.Id) (*db.User, error) {
	user, err := p.UserRepository.GetById(userId)
	if err != nil {
		return nil, fmt.Errorf("%w: get user: %v", ErrFailedGetProfile, err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// UpdateProfile updates a user's profile.
func (p *Profile) UpdateProfile(userId db.Id, req *UpdateProfileRequest) (*db.User, error) {
	user, err := p.UserRepository.GetById(userId)
	if err != nil {
		return nil, fmt.Errorf("%w: get user: %v", ErrFailedUpdateProfile, err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if (lo.IsNotEmpty(req.Email) && req.Email != user.Email) || lo.IsNotEmpty(req.NewPassword) {
		if lo.IsEmpty(req.CurrentPassword) {
			return nil, ErrCurrentPasswordRequired
		}
		if !util.ComparePassword(user.Password, req.CurrentPassword) {
			return nil, ErrCurrentPasswordIncorrect
		}
	}

	user.Name = lo.Ternary(lo.IsNotEmpty(req.Name), req.Name, user.Name)
	email := lo.Ternary(lo.IsNotEmpty(req.Email), req.Email, user.Email)

	if email != user.Email {
		existingUser, err := p.UserRepository.GetByEmail(email)
		if err != nil {
			return nil, fmt.Errorf("%w: check existing user: %v", ErrFailedUpdateProfile, err)
		}
		if existingUser != nil && existingUser.Id != user.Id {
			return nil, ErrUserEmailAlreadyExists
		}
	}

	user.Email = email
	user.Language = req.Language
	user.Theme = req.Theme
	if lo.IsNotEmpty(req.NewPassword) {
		hashedPassword, err := util.HashPassword(req.NewPassword)
		if err != nil {
			return nil, fmt.Errorf("%w: hash password: %v", ErrFailedUpdateProfile, err)
		}
		user.Password = hashedPassword
	}

	err = p.UserRepository.Update(user)
	if err != nil {
		return nil, fmt.Errorf("%w: update user: %v", ErrFailedUpdateProfile, err)
	}

	return user, nil
}

// ResendVerifyEmail sends another verification email, limited to 5 per day.
func (p *Profile) ResendVerifyEmail(userId db.Id) error {
	user, err := p.UserRepository.GetById(userId)
	if err != nil {
		return fmt.Errorf("%w: get user: %v", ErrFailedResendVerifyEmail, err)
	}
	if user == nil {
		return ErrUserNotFound
	}
	if user.IsEmailVerified {
		return ErrEmailAlreadyVerified
	}

	today := time.Now().UTC().Format("2006-01-02")
	count := 0
	meta, err := p.UserMetaRepository.Get(userId, VerifyEmailResendMetaKey)
	if err != nil {
		return fmt.Errorf("%w: get meta: %v", ErrFailedResendVerifyEmail, err)
	}
	if meta != nil {
		var state struct {
			Day   string `json:"day"`
			Count int    `json:"count"`
		}
		_ = json.Unmarshal([]byte(meta.Value), &state)
		if state.Day == today {
			count = state.Count
		}
	}
	if count >= VerifyEmailResendDailyLimit {
		return ErrVerifyEmailResendLimitReached
	}

	token := uuid.NewString()
	user.EmailVerifyToken = &token
	err = p.UserRepository.Update(user)
	if err != nil {
		return fmt.Errorf("%w: update user: %v", ErrFailedResendVerifyEmail, err)
	}

	raw, err := json.Marshal(map[string]any{
		"day":   today,
		"count": count + 1,
	})
	if err != nil {
		return fmt.Errorf("%w: encode meta: %v", ErrFailedResendVerifyEmail, err)
	}
	err = p.UserMetaRepository.Upsert(userId, VerifyEmailResendMetaKey, string(raw))
	if err != nil {
		return fmt.Errorf("%w: upsert meta: %v", ErrFailedResendVerifyEmail, err)
	}

	mailer := resend.NewMailer()
	err = mailer.SendVerifyEmail(
		user.Email,
		user.Name,
		fmt.Sprintf("%s/verify-email/%s", viper.GetString("app.url"), token),
		viper.GetString("app.name"),
	)
	if err != nil {
		log.Warn().
			Err(err).
			Str("userId", user.Id.String()).
			Str("email", user.Email).
			Msg("Failed to resend verify email")
		return fmt.Errorf("%w: %v", ErrFailedResendVerifyEmail, err)
	}

	return nil
}
