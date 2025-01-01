// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/clivern/actx0/mcp"

	"github.com/spf13/cobra"
)

var (
	mcpMgt     string
	mcpPort    int
	mcpPath    string
	mcpTimeout int
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start the actx0 MCP server",
	Long:  `Start the actx0 Model Context Protocol server as a proxy to the management API.`,
	Run: func(_ *cobra.Command, _ []string) {
		err := mcp.Run(context.Background(), mcp.Options{
			Name:    "actx0",
			Version: Version,
			Mgt:     mcpMgt,
			Port:    mcpPort,
			Path:    mcpPath,
			Timeout: time.Duration(mcpTimeout) * time.Second,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "MCP server error: %s\n", err.Error())
			os.Exit(1)
		}
	},
}

// init registers the mcp subcommand.
func init() {
	mcpCmd.Flags().StringVar(
		&mcpMgt,
		"mgt",
		"http://127.0.0.1:8000",
		"Actx0 management server URL",
	)
	mcpCmd.Flags().IntVar(
		&mcpPort,
		"port",
		8001,
		"MCP HTTP listen port",
	)
	mcpCmd.Flags().StringVar(
		&mcpPath,
		"path",
		"/mcp",
		"MCP HTTP path",
	)
	mcpCmd.Flags().IntVar(
		&mcpTimeout,
		"timeout",
		50,
		"Request timeout in seconds",
	)
	rootCmd.AddCommand(mcpCmd)
}
