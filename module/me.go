// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"encoding/json"
	"time"

	"github.com/clivern/actx0/db"
)

// Me resolves API key and access key principals for /me.
type Me struct {
	APIKeyRepository             db.APIKeyRepository
	UserRepository               db.UserRepository
	WorkspaceAccessKeyRepository db.WorkspaceAccessKeyRepository
}

// MeAPIKeyOwner is the user who owns an API key.
type MeAPIKeyOwner struct {
	Id    db.Id  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

// MeAPIKeyInfo is API key metadata returned by /me.
type MeAPIKeyInfo struct {
	Id        db.Id   `json:"id"`
	Name      string  `json:"name"`
	ExpiresAt *string `json:"expiresAt,omitempty"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

// MeAPIKeyResponse is returned when /me is called with X-API-Key.
type MeAPIKeyResponse struct {
	PrincipalType string        `json:"principalType"`
	Owner         MeAPIKeyOwner `json:"owner"`
	APIKey        MeAPIKeyInfo  `json:"apiKey"`
}

// MeAccessKeyResponse is returned when /me is called with X-Access-Key.
type MeAccessKeyResponse struct {
	PrincipalType string            `json:"principalType"`
	AccessKey     AccessKeyResponse `json:"accessKey"`
}

// NewMe creates a /me module with the given repositories.
func NewMe(akr db.APIKeyRepository, ur db.UserRepository, war db.WorkspaceAccessKeyRepository) *Me {
	return &Me{
		APIKeyRepository:             akr,
		UserRepository:               ur,
		WorkspaceAccessKeyRepository: war,
	}
}

// GetByAPIKey returns the API key owner and key metadata.
func (m *Me) GetByAPIKey(key string) (*MeAPIKeyResponse, error) {
	apiKey, err := m.APIKeyRepository.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if apiKey == nil {
		return nil, ErrAPIKeyNotFound
	}

	owner, err := m.UserRepository.GetById(apiKey.UserId)
	if err != nil {
		return nil, err
	}
	if owner == nil {
		return nil, ErrAPIKeyNotFound
	}

	info := MeAPIKeyInfo{
		Id:        apiKey.Id,
		Name:      apiKey.Name,
		CreatedAt: apiKey.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: apiKey.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if apiKey.ExpiresAt != nil {
		formatted := apiKey.ExpiresAt.UTC().Format(time.RFC3339)
		info.ExpiresAt = &formatted
	}

	return &MeAPIKeyResponse{
		PrincipalType: "api_key",
		Owner: MeAPIKeyOwner{
			Id:    owner.Id,
			Email: owner.Email,
			Name:  owner.Name,
			Role:  owner.Role,
		},
		APIKey: info,
	}, nil
}

// GetByAccessKey returns workspace access key metadata and permissions.
func (m *Me) GetByAccessKey(key string) (*MeAccessKeyResponse, error) {
	accessKey, err := m.WorkspaceAccessKeyRepository.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if accessKey == nil {
		return nil, ErrAccessKeyNotFound
	}

	var meta AccessKeyMeta
	err = json.Unmarshal([]byte(*accessKey.Meta), &meta)
	if err != nil {
		return nil, err
	}

	resp := AccessKeyResponse{
		Id:          accessKey.Id,
		WorkspaceId: accessKey.WorkspaceId,
		Name:        accessKey.Name,
		Permissions: meta.Permissions,
		CreatedAt:   accessKey.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   accessKey.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if accessKey.ExpiresAt != nil {
		formatted := accessKey.ExpiresAt.UTC().Format(time.RFC3339)
		resp.ExpiresAt = &formatted
	}

	return &MeAccessKeyResponse{
		PrincipalType: "access_key",
		AccessKey:     resp,
	}, nil
}
