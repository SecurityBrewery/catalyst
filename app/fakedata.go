package app

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"

	"github.com/SecurityBrewery/catalyst/fakedata"
)

func fakeDataCmd(app core.App) *cobra.Command {
	var userCount, ticketCount int

	cmd := &cobra.Command{
		Use: "fake-data",
		Run: func(_ *cobra.Command, _ []string) {
			if err := fakedata.Generate(app, userCount, ticketCount); err != nil {
				app.Logger().Error(err.Error())
			}
		},
	}

	cmd.PersistentFlags().IntVar(&userCount, "users", 10, "Number of users to generate")

	cmd.PersistentFlags().IntVar(&ticketCount, "tickets", 100, "Number of tickets to generate")

	return cmd
}
