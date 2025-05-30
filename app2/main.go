package app2

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/SecurityBrewery/catalyst/app2/database"
	"github.com/SecurityBrewery/catalyst/app2/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app2/openapi"
)

func App(filename string, _ bool) (*App2, error) {
	queries, _, err := database.DB(filepath.Join(filename, "data.db"))
	if err != nil {
		return nil, err
	}

	service := &Service{
		Queries: queries,
	}

	mux := http.NewServeMux()

	apiHandler := openapi.Handler(openapi.NewStrictHandler(service, nil))
	mux.Handle("/api/", http.StripPrefix("/api", apiHandler))

	return &App2{
		Queries: queries,
	}, nil
}

type App2 struct {
	Queries *sqlc.Queries
}

func (a *App2) Start() error {
	ctx := context.Background()

	tickets, err := a.Queries.ListTickets(ctx, sqlc.ListTicketsParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		return err
	}

	for _, ticket := range tickets {
		// Process each ticket as needed
		// For example, you could print the ticket ID
		println("Ticket ID:", ticket.ID)
	}

	return nil
}
