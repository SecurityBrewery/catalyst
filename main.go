package main

import (
	"log"
	"os"

	"github.com/pocketbase/pocketbase"

	"github.com/SecurityBrewery/catalyst/ff"
	"github.com/SecurityBrewery/catalyst/migrations"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	migrations.Register()

	dir := "./catalyst_data"
	if os.Getenv("CATALYST_DIR") != "" {
		dir = os.Getenv("CATALYST_DIR")
	}

	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDev:     ff.HasDevFlag(),
		DefaultDataDir: dir,
	})

	app.OnBeforeServe().Add(addRoutes())

	attachWebhooks(app)

	if !ff.HasDevFlag() {
		os.Args = []string{"catalyst", "serve", "--http", "0.0.0.0:8088"}
	}

	return app.Start()
}
