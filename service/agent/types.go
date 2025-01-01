// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package agent

import "github.com/clivern/actx0/db"

const (
	MemoryKindSummary    = "summary"
	MemoryKindFact       = "fact"
	MemoryKindPreference = "preference"
	MemoryKindShortLived = "short_lived"
	MemoryKindLongLived  = "long_lived"
)

// StoreMemoryRequest persists a session memory item.
type StoreMemoryRequest struct {
	ExternalId string
	AgentId    db.Id
	Kind       string
	Content    string
	Meta       *string
}

// ConsolidateMemoryRequest triggers async session memory consolidation.
type ConsolidateMemoryRequest struct {
	AgentId    db.Id
	ExternalId string
}

// RecallMemoryRequest finds memories relevant to a query.
type RecallMemoryRequest struct {
	ExternalId string
	AgentId    db.Id
	Query      string
	Kinds      []string
	Limit      uint64
}

// MemoryHit is a scored memory item.
type MemoryHit struct {
	MemoryId db.Id
	Kind     string
	Content  string
	Score    float32
}

// SearchSessionRequest finds session vectors relevant to a query.
type SearchSessionRequest struct {
	WorkspaceId db.Id
	AgentId     db.Id
	SessionId   db.Id
	Query       string
	Limit       uint64
}

// MessageHit is a scored session message.
type MessageHit struct {
	MessageId db.Id
	Role      string
	Score     float32
	Text      string
}
