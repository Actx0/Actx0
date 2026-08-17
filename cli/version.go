// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package cli

import (
	"embed"
	"fmt"

	"github.com/clivern/actx0/build"

	"github.com/spf13/cobra"
)

var (
	// Static embedded files
	Static embed.FS
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf(
			"Current Actx0 version %v commit %v, built @%v by %v.\n",
			build.Version,
			build.Commit,
			build.Date,
			build.BuiltBy,
		)
	},
}

// init registers the version subcommand.
func init() {
	rootCmd.AddCommand(versionCmd)
}
