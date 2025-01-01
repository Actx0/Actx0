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
	"github.com/clivern/actx0/db"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database shell commands",
	Long:  `Database shell and utility commands`,
}

var dbShellCmd = &cobra.Command{
	Use:   "shell [-- psql-args...]",
	Short: "Open an interactive database shell",
	Long:  `Open an interactive PostgreSQL shell (psql) connected to the configured database. Extra args after -- are passed to psql.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := loadDBConfig(cmd)
		execPsql(cfg, args)
	},
}

var dbQueryCmd = &cobra.Command{
	Use:   "query [sql]",
	Short: "Run a SQL statement and exit",
	Long:  `Run one or more SQL statements via psql (-c) against the configured database`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := loadDBConfig(cmd)
		execPsql(cfg, []string{"-v", "ON_ERROR_STOP=1", "-c", strings.Join(args, " ")})
	},
}

func loadDBConfig(cmd *cobra.Command) db.DatabaseConfig {
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

	return core.ReadWriteDatabase()
}

func execPsql(cfg db.DatabaseConfig, extraArgs []string) {
	psqlPath, err := exec.LookPath("psql")
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("psql command not found. Please install PostgreSQL client tools")
	}

	psqlArgs := []string{
		"psql",
		fmt.Sprintf("--host=%s", cfg.Host),
		fmt.Sprintf("--port=%d", cfg.Port),
		fmt.Sprintf("--username=%s", cfg.Username),
		fmt.Sprintf("--dbname=%s", cfg.Database),
	}
	psqlArgs = append(psqlArgs, extraArgs...)

	env := append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", cfg.Password))

	log.Info().
		Str("host", cfg.Host).
		Int("port", cfg.Port).
		Str("database", cfg.Database).
		Str("username", cfg.Username).
		Msg("Connecting to database")

	err = syscall.Exec(psqlPath, psqlArgs, env)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to execute psql")
	}
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
