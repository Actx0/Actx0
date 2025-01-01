// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package mcp

import (
	"context"
	"net/http"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	HeaderWorkspaceID = "X-Workspace-Id"
	HeaderAgentID     = "X-Agent-Id"
	HeaderSessionID   = "X-Session-Id"
	HeaderAccessKey   = "X-Access-Key"
)

// Scope is the connection-level identity for MCP tool calls.
type Scope struct {
	WorkspaceID string
	AgentID     string
	SessionID   string
	AccessKey   string
}

type scopeKey struct{}

// ScopeFromHeader builds a Scope from request headers.
func ScopeFromHeader(h http.Header) Scope {
	return Scope{
		WorkspaceID: h.Get(HeaderWorkspaceID),
		AgentID:     h.Get(HeaderAgentID),
		SessionID:   h.Get(HeaderSessionID),
		AccessKey:   h.Get(HeaderAccessKey),
	}
}

// WithScope attaches Scope to a context.
func WithScope(ctx context.Context, scope Scope) context.Context {
	return context.WithValue(ctx, scopeKey{}, scope)
}

// ScopeFrom reads Scope previously attached by middleware.
func ScopeFrom(ctx context.Context) Scope {
	return ctx.Value(scopeKey{}).(Scope)
}

// HTTPScope middleware stores workspace/agent/session headers on the request context.
func HTTPScope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := WithScope(r.Context(), ScopeFromHeader(r.Header))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ScopeMiddleware copies HTTP headers from CallToolRequest Extra onto ctx for every tool.
func ScopeMiddleware(next mcpsdk.MethodHandler) mcpsdk.MethodHandler {
	return func(ctx context.Context, method string, req mcpsdk.Request) (mcpsdk.Result, error) {
		ctx = WithScope(ctx, ScopeFromHeader(req.GetExtra().Header))

		return next(ctx, method, req)
	}
}
