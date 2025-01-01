// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package mcp

import "encoding/json"

type Session struct {
	Id string `json:"id"`
}

type Memory struct {
	Id        string          `json:"id"`
	SessionId string          `json:"sessionId"`
	Kind      string          `json:"kind"`
	Content   string          `json:"content"`
	Meta      json.RawMessage `json:"meta,omitempty"`
	CreatedAt string          `json:"createdAt"`
	UpdatedAt string          `json:"updatedAt"`
}

type ListMemories struct {
	Memories []Memory `json:"memories"`
	Meta     struct {
		Limit  int   `json:"limit"`
		Offset int   `json:"offset"`
		Total  int64 `json:"total"`
	} `json:"_meta"`
}

type SearchHit struct {
	Id    string  `json:"id"`
	Kind  string  `json:"kind"`
	Score float32 `json:"score"`
	Text  string  `json:"text"`
}

type SearchMemories struct {
	Results []SearchHit `json:"results"`
}

type MemoryBody struct {
	Kind    string  `json:"kind,omitempty"`
	Content string  `json:"content"`
	Meta    *string `json:"meta,omitempty"`
}

type SearchBody struct {
	Query string `json:"query"`
	Limit int    `json:"limit,omitempty"`
}

type AddMemoryInput struct {
	Kind     string         `json:"kind,omitempty" jsonschema:"Memory kind: summary, fact, preference, short_lived, or long_lived"`
	Content  string         `json:"content" jsonschema:"required,Text content to store"`
	Metadata map[string]any `json:"metadata,omitempty" jsonschema:"Optional metadata"`
}

type AddMemoryOutput struct {
	MemoryId string `json:"memory_id"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

type SearchMemoriesInput struct {
	Query string `json:"query" jsonschema:"required,Semantic search query"`
	Limit int    `json:"limit,omitempty" jsonschema:"Maximum number of results"`
}

type MemoryRecord struct {
	MemoryId  string         `json:"memory_id"`
	Kind      string         `json:"kind,omitempty"`
	Memory    string         `json:"memory"`
	Score     float32        `json:"score,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	CreatedAt string         `json:"created_at,omitempty"`
	UpdatedAt string         `json:"updated_at,omitempty"`
}

type SearchMemoriesOutput struct {
	Results []MemoryRecord `json:"results"`
}

type GetMemoriesInput struct {
	Limit  int `json:"limit,omitempty" jsonschema:"Page size"`
	Offset int `json:"offset,omitempty" jsonschema:"Page offset"`
}

type GetMemoriesOutput struct {
	Memories []MemoryRecord `json:"memories"`
	Total    int64          `json:"total"`
	Limit    int            `json:"limit"`
	Offset   int            `json:"offset"`
}

type GetMemoryInput struct {
	MemoryId string `json:"memory_id" jsonschema:"required,Memory identifier"`
}

type UpdateMemoryInput struct {
	MemoryId string         `json:"memory_id" jsonschema:"required,Memory identifier"`
	Kind     string         `json:"kind,omitempty" jsonschema:"Memory kind: summary, fact, preference, short_lived, or long_lived"`
	Content  string         `json:"content" jsonschema:"required,Updated memory text"`
	Metadata map[string]any `json:"metadata,omitempty" jsonschema:"Optional metadata"`
}

type UpdateMemoryOutput struct {
	MemoryId string `json:"memory_id"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

type DeleteMemoryInput struct {
	MemoryId string `json:"memory_id" jsonschema:"required,Memory identifier"`
}

type DeleteMemoryOutput struct {
	MemoryId string `json:"memory_id"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

type DeleteSessionOutput struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
