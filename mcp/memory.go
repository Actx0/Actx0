// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package mcp

import (
	"context"
	"encoding/json"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func RegisterMemoryTools(server *mcpsdk.Server, api *APIClient) {
	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "add_memory",
		Description: "Save a memory for the scoped session",
	}, func(ctx context.Context, _ *mcpsdk.CallToolRequest, input AddMemoryInput) (*mcpsdk.CallToolResult, AddMemoryOutput, error) {
		scope := ScopeFrom(ctx)

		kind := input.Kind
		if kind == "" {
			kind = "long_lived"
		}

		meta, err := EncodeMeta(input.Metadata)
		if err != nil {
			return nil, AddMemoryOutput{}, err
		}

		memory, err := api.CreateMemory(ctx, scope, kind, input.Content, meta)
		if err != nil {
			return nil, AddMemoryOutput{}, err
		}

		return nil, AddMemoryOutput{
			MemoryId: memory.Id,
			Status:   "completed",
			Message:  "Memory saved",
		}, nil
	})

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "search_memories",
		Description: "Semantic search across memories in the scoped session",
	}, func(ctx context.Context, _ *mcpsdk.CallToolRequest, input SearchMemoriesInput) (*mcpsdk.CallToolResult, SearchMemoriesOutput, error) {
		scope := ScopeFrom(ctx)

		result, err := api.SearchMemories(ctx, scope, input.Query, input.Limit)
		if err != nil {
			return nil, SearchMemoriesOutput{}, err
		}

		records := make([]MemoryRecord, 0, len(result.Results))
		for _, hit := range result.Results {
			records = append(records, MemoryRecord{
				MemoryId: hit.Id,
				Kind:     hit.Kind,
				Memory:   hit.Text,
				Score:    hit.Score,
			})
		}

		return nil, SearchMemoriesOutput{Results: records}, nil
	})

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "get_memories",
		Description: "List memories in the scoped session",
	}, func(ctx context.Context, _ *mcpsdk.CallToolRequest, input GetMemoriesInput) (*mcpsdk.CallToolResult, GetMemoriesOutput, error) {
		scope := ScopeFrom(ctx)

		list, err := api.ListMemories(ctx, scope, input.Limit, input.Offset)
		if err != nil {
			return nil, GetMemoriesOutput{}, err
		}

		records := make([]MemoryRecord, 0, len(list.Memories))
		for _, item := range list.Memories {
			record := MemoryRecord{
				MemoryId:  item.Id,
				Kind:      item.Kind,
				Memory:    item.Content,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			}
			if len(item.Meta) > 0 {
				_ = json.Unmarshal(item.Meta, &record.Metadata)
			}
			records = append(records, record)
		}

		return nil, GetMemoriesOutput{
			Memories: records,
			Total:    list.Meta.Total,
			Limit:    list.Meta.Limit,
			Offset:   list.Meta.Offset,
		}, nil
	})

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "get_memory",
		Description: "Retrieve one memory by its memory_id",
	}, func(ctx context.Context, _ *mcpsdk.CallToolRequest, input GetMemoryInput) (*mcpsdk.CallToolResult, MemoryRecord, error) {
		scope := ScopeFrom(ctx)

		memory, err := api.GetMemory(ctx, scope, input.MemoryId)
		if err != nil {
			return nil, MemoryRecord{}, err
		}

		record := MemoryRecord{
			MemoryId:  memory.Id,
			Kind:      memory.Kind,
			Memory:    memory.Content,
			CreatedAt: memory.CreatedAt,
			UpdatedAt: memory.UpdatedAt,
		}
		if len(memory.Meta) > 0 {
			_ = json.Unmarshal(memory.Meta, &record.Metadata)
		}

		return nil, record, nil
	})

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "update_memory",
		Description: "Overwrite a memory's text after confirming the ID",
	}, func(ctx context.Context, _ *mcpsdk.CallToolRequest, input UpdateMemoryInput) (*mcpsdk.CallToolResult, UpdateMemoryOutput, error) {
		scope := ScopeFrom(ctx)

		meta, err := EncodeMeta(input.Metadata)
		if err != nil {
			return nil, UpdateMemoryOutput{}, err
		}

		memory, err := api.UpdateMemory(ctx, scope, input.MemoryId, input.Kind, input.Content, meta)
		if err != nil {
			return nil, UpdateMemoryOutput{}, err
		}

		return nil, UpdateMemoryOutput{
			MemoryId: memory.Id,
			Status:   "completed",
			Message:  "Memory updated",
		}, nil
	})

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "delete_memory",
		Description: "Delete a single memory by memory_id",
	}, func(ctx context.Context, _ *mcpsdk.CallToolRequest, input DeleteMemoryInput) (*mcpsdk.CallToolResult, DeleteMemoryOutput, error) {
		scope := ScopeFrom(ctx)

		err := api.DeleteMemory(ctx, scope, input.MemoryId)
		if err != nil {
			return nil, DeleteMemoryOutput{}, err
		}

		return nil, DeleteMemoryOutput{
			MemoryId: input.MemoryId,
			Status:   "completed",
			Message:  "Memory deleted",
		}, nil
	})

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "delete_session",
		Description: "Delete the scoped session and all of its memories",
	}, func(ctx context.Context, _ *mcpsdk.CallToolRequest, _ struct{}) (*mcpsdk.CallToolResult, DeleteSessionOutput, error) {
		scope := ScopeFrom(ctx)

		err := api.DeleteSession(ctx, scope)
		if err != nil {
			return nil, DeleteSessionOutput{}, err
		}

		return nil, DeleteSessionOutput{
			Status:  "completed",
			Message: "Session deleted",
		}, nil
	})
}
