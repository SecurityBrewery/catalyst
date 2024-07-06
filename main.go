package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

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

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{})

	app.RootCmd.AddCommand(fakeDataCmd(app))
	app.RootCmd.AddCommand(setFeatureFlagsCmd(app))

	app.OnBeforeServe().Add(addRoutes())

	attachWebhooks(app)

	return app.Start()
}
