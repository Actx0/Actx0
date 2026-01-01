// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package cli

import (
	"github.com/clivern/actx0/core"
	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/migration"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
	Long:  `Manage database migrations`,
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all pending migrations",
	Run: func(cmd *cobra.Command, _ []string) {
		configFile, _ := cmd.Flags().GetString("config")

		if err := core.Load(configFile); err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to load configuration")
		}

		if err := core.SetupLogging(); err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to setup logging")
		}

		dbConfig := db.Config{
			Driver:          viper.GetString("app.database.driver"),
			Host:            viper.GetString("app.database.host"),
			Port:            viper.GetInt("app.database.port"),
			Username:        viper.GetString("app.database.username"),
			Password:        viper.GetString("app.database.password"),
			Database:        viper.GetString("app.database.name"),
			MaxOpenConns:    viper.GetInt("app.database.max_open_conns"),
			MaxIdleConns:    viper.GetInt("app.database.max_idle_conns"),
			ConnMaxLifetime: viper.GetInt("app.database.conn_max_lifetime"),
		}

		conn, err := db.NewConnection(dbConfig)
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to connect to database")
		}
		defer conn.Close()

		mgr := migration.NewManager(conn.DB, conn.Driver)

		for _, m := range migration.GetAll() {
			mgr.Register(m)
		}

		if err := mgr.Up(); err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to run migrations")
		}

		log.Info().Msg("Migration completed successfully")
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Roll back the last migration",
	Run: func(cmd *cobra.Command, _ []string) {
		configFile, _ := cmd.Flags().GetString("config")

		if err := core.Load(configFile); err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to load configuration")
		}

		if err := core.SetupLogging(); err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to setup logging")
		}

		dbConfig := db.Config{
			Driver:          viper.GetString("app.database.driver"),
			Host:            viper.GetString("app.database.host"),
			Port:            viper.GetInt("app.database.port"),
			Username:        viper.GetString("app.database.username"),
			Password:        viper.GetString("app.database.password"),
			Database:        viper.GetString("app.database.name"),
			MaxOpenConns:    viper.GetInt("app.database.max_open_conns"),
			MaxIdleConns:    viper.GetInt("app.database.max_idle_conns"),
			ConnMaxLifetime: viper.GetInt("app.database.conn_max_lifetime"),
		}

		conn, err := db.NewConnection(dbConfig)
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to connect to database")
		}
		defer conn.Close()

		mgr := migration.NewManager(conn.DB, conn.Driver)

		for _, m := range migration.GetAll() {
			mgr.Register(m)
		}

		if err := mgr.Down(); err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to roll back migration")
		}

		log.Info().Msg("Rollback completed successfully")
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Run: func(cmd *cobra.Command, _ []string) {
		configFile, _ := cmd.Flags().GetString("config")

		if err := core.Load(configFile); err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to load configuration")
		}

		if err := core.SetupLogging(); err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to setup logging")
		}

		dbConfig := db.Config{
			Driver:          viper.GetString("app.database.driver"),
			Host:            viper.GetString("app.database.host"),
			Port:            viper.GetInt("app.database.port"),
			Username:        viper.GetString("app.database.username"),
			Password:        viper.GetString("app.database.password"),
			Database:        viper.GetString("app.database.name"),
			MaxOpenConns:    viper.GetInt("app.database.max_open_conns"),
			MaxIdleConns:    viper.GetInt("app.database.max_idle_conns"),
			ConnMaxLifetime: viper.GetInt("app.database.conn_max_lifetime"),
		}

		conn, err := db.NewConnection(dbConfig)
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to connect to database")
		}
		defer conn.Close()

		mgr := migration.NewManager(conn.DB, conn.Driver)

		for _, m := range migration.GetAll() {
			mgr.Register(m)
		}

		if err := mgr.Status(); err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to get migration status")
		}
	},
}

// init registers the migrate subcommands and flags.
func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateStatusCmd)

	migrateUpCmd.Flags().StringVarP(
		&config,
		"config",
		"c",
		"config.prod.yml",
		"Absolute path to config file (required)",
	)
	migrateUpCmd.MarkFlagRequired("config")
	migrateDownCmd.Flags().StringVarP(
		&config,
		"config",
		"c",
		"config.prod.yml",
		"Absolute path to config file (required)",
	)
	migrateDownCmd.MarkFlagRequired("config")
	migrateStatusCmd.Flags().StringVarP(
		&config,
		"config",
		"c",
		"config.prod.yml",
		"Absolute path to config file (required)",
	)
	migrateStatusCmd.MarkFlagRequired("config")
}
