package main

import (
	"log"

	"github.com/pocketbase/pocketbase"

	"github.com/SecurityBrewery/catalyst/migrations"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	migrations.Register()

	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDev:     dev(),
		DefaultDataDir: "catalyst_data",
	})

	attachWebhooks(app)

	// Register additional commands
	app.RootCmd.AddCommand(bootstrapCmd(app))
	app.RootCmd.AddCommand(fakeDataCmd(app))
	app.RootCmd.AddCommand(setFeatureFlagsCmd(app))

	app.OnBeforeServe().Add(addRoutes())

	return app.Start()
}
