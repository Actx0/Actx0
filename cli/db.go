// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/clivern/actx0/core"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database shell commands",
	Long:  `Database shell and utility commands`,
}

var dbShellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Open an interactive database shell",
	Long:  `Open an interactive PostgreSQL shell (psql) connected to the configured database.`,
	Run: func(cmd *cobra.Command, args []string) {
		configFile, _ := cmd.Flags().GetString("config")

		err := core.Load(configFile)
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to load configuration")
		}

		err = core.SetupLogging()
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to setup logging")
		}

		cfg := core.ReadWriteDatabase()

		path, err := exec.LookPath("psql")
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("psql command not found. Please install PostgreSQL client tools")
		}

		pargs := []string{
			"psql",
			"-q",
			fmt.Sprintf("--host=%s", cfg.Host),
			fmt.Sprintf("--port=%d", cfg.Port),
			fmt.Sprintf("--username=%s", cfg.Username),
			fmt.Sprintf("--dbname=%s", cfg.Database),
		}
		pargs = append(pargs, args...)

		env := append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", cfg.Password))

		fmt.Printf(`
    ___   _____________  ___
   / _ | / ___/_  __/ |/_/ /
  / __ |/ /__  / / _>  </_/
 /_/ |_|___/ /_/ /_/|_/_/

 host     %s:%d
 database %s
 user     %s

`, cfg.Host, cfg.Port, cfg.Database, cfg.Username)

		err = syscall.Exec(path, pargs, env)
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to execute psql")
		}
	},
}

var dbQueryCmd = &cobra.Command{
	Use:   "query [sql]",
	Short: "Run a SQL statement and exit",
	Long:  `Run one or more SQL statements via psql against the configured database`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configFile, _ := cmd.Flags().GetString("config")

		err := core.Load(configFile)
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to load configuration")
		}

		err = core.SetupLogging()
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to setup logging")
		}

		cfg := core.ReadWriteDatabase()

		path, err := exec.LookPath("psql")
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("psql command not found. Please install PostgreSQL client tools")
		}

		pargs := []string{
			"psql",
			"-q",
			fmt.Sprintf("--host=%s", cfg.Host),
			fmt.Sprintf("--port=%d", cfg.Port),
			fmt.Sprintf("--username=%s", cfg.Username),
			fmt.Sprintf("--dbname=%s", cfg.Database),
			"-v", "ON_ERROR_STOP=1",
			"-c", strings.Join(args, " "),
		}

		env := append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", cfg.Password))

		err = syscall.Exec(path, pargs, env)
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to execute psql")
		}
	},
}

// init registers the db subcommands and flags.
func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(dbShellCmd)
	dbCmd.AddCommand(dbQueryCmd)

	dbCmd.PersistentFlags().StringP(
		"config",
		"c",
		"config.dist.yml",
		"Path to config file",
	)
}
