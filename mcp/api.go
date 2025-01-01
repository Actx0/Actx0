// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package mcp

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// APIClient calls the actx0 management HTTP API.
type APIClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewAPIClient builds a client pointed at the management server.
func NewAPIClient(mgt string, timeout time.Duration) *APIClient {
	return &APIClient{
		baseURL: strings.TrimRight(mgt, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// DeleteMemory deletes a session memory via the actx0 API.
func (c *APIClient) DeleteMemory(ctx context.Context, scope Scope, memoryID string) error {
	url := fmt.Sprintf(
		"%s/api/v1/workspaces/%s/agents/%s/sessions/%s/memories/%s",
		c.baseURL,
		scope.WorkspaceID,
		scope.AgentID,
		scope.SessionID,
		memoryID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set(HeaderAccessKey, scope.AccessKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete memory: unexpected status %d", resp.StatusCode)
	}

	return nil
}
