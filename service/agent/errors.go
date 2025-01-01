// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package agent

import "errors"

var (
	ErrNotImplemented    = errors.New("not implemented")
	ErrAgentNotFound     = errors.New("agent not found")
	ErrSessionNotFound   = errors.New("agent session not found")
	ErrMemoryNotFound    = errors.New("session memory not found")
	ErrMessageNotFound   = errors.New("session message not found")
	ErrStoreFailed       = errors.New("session memory store failed")
	ErrRecallFailed      = errors.New("session memory recall failed")
	ErrDeleteFailed      = errors.New("agent delete failed")
	ErrSearchFailed      = errors.New("agent search failed")
	ErrIndexFailed       = errors.New("agent index failed")
	ErrGenerateFailed    = errors.New("session memory generate failed")
	ErrConsolidateFailed = errors.New("session memory consolidate failed")
	ErrRegenerateFailed  = errors.New("session memory regenerate failed")
)
