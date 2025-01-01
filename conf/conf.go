// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package conf

import "time"

const (
	// SlowRequestThreshold is the latency at or above which a request is counted as slow.
	SlowRequestThreshold = time.Second

	// DefaultSessionDuration is the auth session TTL without remember-me.
	DefaultSessionDuration = 7 * 24 * time.Hour

	// RememberMeSessionDuration is the auth session TTL with remember-me.
	RememberMeSessionDuration = 30 * 24 * time.Hour

	// PasswordResetTokenTTL is how long a password reset token stays valid.
	PasswordResetTokenTTL = 24 * time.Hour

	// InviteExpiry is how long a workspace invite stays valid.
	InviteExpiry = 7 * 24 * time.Hour

	// AsyncClaimInterval is how often the async manager polls for due tasks.
	AsyncClaimInterval = 10 * time.Second

	// AsyncTaskMaxRetries is how many times a failed task is retried.
	AsyncTaskMaxRetries = 5

	// AsyncCompletedTasksRetention is how far back completed tasks are kept before cleanup.
	// Negative duration: cutoff = now + retention.
	AsyncCompletedTasksRetention = -7 * 24 * time.Hour

	// MaxUploadBytes is the maximum multipart upload size.
	MaxUploadBytes = 2 * 1024 * 1024

	// DefaultSearchLimit is used when a knowledge search request omits limit.
	DefaultSearchLimit = 10

	// GenerateMemorySystemPrompt instructs the lite LLM when extracting session memory.
	GenerateMemorySystemPrompt = `You extract durable memory from one conversation message (prefixed with role=user or role=assistant) for an AI agent memory store.

Rules:
- Content must be minimal: one short third-person sentence, no filler, no restating the conversation. Keep exact names, numbers, dates, and identifiers.
- long_lived: stable identity, preferences, lasting facts, relationships.
- short_lived: ephemeral task state, temporary goals, context that fades after this conversation.
- Skip greetings, chit-chat, generic instructions, and anything with no value later.
- role=assistant: skip unless it states a concrete fact, decision, commitment, or result worth recalling. Never store explanations, generic answers, or restated user input.

Your entire reply must be one JSON object. The first character must be { and the last must be }.
No markdown, no explanation, no other text.

Schema:
{"skip":true,"kind":"","content":""}
or
{"skip":false,"kind":"short_lived","content":"minimal memory sentence"}
or
{"skip":false,"kind":"long_lived","content":"minimal memory sentence"}`

	// ConsolidateMemorySystemPrompt instructs the lite LLM when consolidating session memories.
	ConsolidateMemorySystemPrompt = `You consolidate an AI agent memory store. Input is a JSON array of memories: [{"id":"uuid","kind":"short_lived|long_lived","content":"..."}].

Return the operations to apply:
- {"op":"delete","id":"uuid"} to remove a duplicate, contradicted, obsolete, or worthless memory.
- {"op":"update","id":"uuid","kind":"short_lived|long_lived","content":"..."} to rewrite a memory with corrected or merged content, or to change its kind.
- {"op":"insert","kind":"short_lived|long_lived","content":"..."} to add a memory that merges several others (delete the ones it replaces).

Rules:
- Merge overlapping memories into one; keep the fewest memories that preserve all meaning.
- On contradictions, keep the newer fact and delete the outdated one.
- Content must be minimal: one short third-person sentence, no filler. Keep exact names, numbers, dates, and identifiers.
- Promote stable facts to long_lived; demote or delete stale short_lived items.
- Do not touch memories that need no change. If nothing needs changing, return {"ops":[]}.

Your entire reply must be one JSON object. The first character must be { and the last must be }.
No markdown, no explanation, no other text.

Schema:
{"ops":[{"op":"delete","id":"uuid"},{"op":"update","id":"uuid","kind":"long_lived","content":"minimal memory sentence"},{"op":"insert","kind":"short_lived","content":"minimal memory sentence"}]}`
)
