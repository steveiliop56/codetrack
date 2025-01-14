package cmd

import (
	"codetrack/internal/api"
	"codetrack/internal/db"
	"codetrack/internal/queries"
	"codetrack/internal/services"
	"codetrack/internal/utils"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "codetrack",
	Short: "Simple, open-source, self-hosted analytics for your code habits",
	Long: `Codetrack is a simple, open-source and self-hosted analytics tool for your code habits. It collects data from your favorite code editors and presents them in a nice dashboard.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Setup logger
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		log.Info().Msg("Starting codetrack")

		// Create injector
		injector := do.New()

		// Create providers
		do.Provide(injector, db.NewDatabase)
		do.Provide(injector, api.NewApi)
		do.Provide(injector, queries.NewQueries)
		do.Provide(injector, services.NewServices)

		// Invoke
		database, databaseErr := do.Invoke[*db.Database](injector)
		utils.HandleError(databaseErr, true)

		api, apiErr := do.Invoke[*api.Api](injector)
		utils.HandleError(apiErr, true)

		// Database
		log.Info().Msg("Initializing database")
		database.Initialize("codetrack.db")

		// Migrate
		log.Info().Msg("Migrating database")
		migrateErr := database.Migrate()
		utils.HandleError(migrateErr, true)

		// Start API
		log.Info().Msg("Starting API")
		api.Initialize("secret")
		api.SetupRoutes()
		api.Start()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}


