package database

import (
	"context"
	"fmt"
	"log"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"

	"github.com/SecurityBrewery/catalyst/bus"
	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/database/migrations"
	"github.com/SecurityBrewery/catalyst/hooks"
	"github.com/SecurityBrewery/catalyst/index"
)

const (
	Name                     = "catalyst"
	TicketCollectionName     = "tickets"
	TemplateCollectionName   = "templates"
	PlaybookCollectionName   = "playbooks"
	AutomationCollectionName = "automations"
	UserDataCollectionName   = "userdata"
	UserCollectionName       = "users"
	TicketTypeCollectionName = "tickettypes"
	JobCollectionName        = "jobs"

	TicketArtifactsGraphName     = "Graph"
	RelatedTicketsCollectionName = "related"
)

type Database struct {
	*busdb.BusDatabase
	Index *index.Index
	bus   *bus.Bus
	Hooks *hooks.Hooks

	templateCollection   *busdb.Collection
	ticketCollection     *busdb.Collection
	playbookCollection   *busdb.Collection
	automationCollection *busdb.Collection
	userdataCollection   *busdb.Collection
	userCollection       *busdb.Collection
	tickettypeCollection *busdb.Collection
	jobCollection        *busdb.Collection

	relatedCollection  *busdb.Collection
	containsCollection *busdb.Collection
}

type Config struct {
	Host     string
	User     string
	Password string
	Name     string
}

func New(ctx context.Context, index *index.Index, bus *bus.Bus, hooks *hooks.Hooks, config *Config) (*Database, error) {
	name := config.Name
	if config.Name == "" {
		name = Name
	}

	conn, err := http.NewConnection(http.ConnectionConfig{Endpoints: []string{config.Host}})
	if err != nil {
		return nil, err
	}
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(config.User, config.Password),
	})
	if err != nil {
		return nil, err
	}

	hooks.DatabaseAfterConnect(ctx, client, name)

	db, err := setupDB(ctx, client, name)
	if err != nil {
		return nil, fmt.Errorf("DB setup failed: %w", err)
	}

	if err = migrations.PerformMigrations(ctx, db); err != nil {
		return nil, fmt.Errorf("migrations failed: %w", err)
	}

	ticketCollection, err := db.Collection(ctx, TicketCollectionName)
	if err != nil {
		return nil, err
	}
	templateCollection, err := db.Collection(ctx, TemplateCollectionName)
	if err != nil {
		return nil, err
	}
	playbookCollection, err := db.Collection(ctx, PlaybookCollectionName)
	if err != nil {
		return nil, err
	}
	relatedCollection, err := db.Collection(ctx, RelatedTicketsCollectionName)
	if err != nil {
		return nil, err
	}
	automationCollection, err := db.Collection(ctx, AutomationCollectionName)
	if err != nil {
		return nil, err
	}
	userdataCollection, err := db.Collection(ctx, UserDataCollectionName)
	if err != nil {
		return nil, err
	}
	userCollection, err := db.Collection(ctx, UserCollectionName)
	if err != nil {
		return nil, err
	}
	tickettypeCollection, err := db.Collection(ctx, TicketTypeCollectionName)
	if err != nil {
		return nil, err
	}
	jobCollection, err := db.Collection(ctx, JobCollectionName)
	if err != nil {
		return nil, err
	}

	hookedDB, err := busdb.NewDatabase(ctx, db, bus)
	if err != nil {
		return nil, err
	}

	return &Database{
		BusDatabase:          hookedDB,
		bus:                  bus,
		Index:                index,
		Hooks:                hooks,
		templateCollection:   busdb.NewCollection(templateCollection, hookedDB),
		ticketCollection:     busdb.NewCollection(ticketCollection, hookedDB),
		playbookCollection:   busdb.NewCollection(playbookCollection, hookedDB),
		automationCollection: busdb.NewCollection(automationCollection, hookedDB),
		relatedCollection:    busdb.NewCollection(relatedCollection, hookedDB),
		userdataCollection:   busdb.NewCollection(userdataCollection, hookedDB),
		userCollection:       busdb.NewCollection(userCollection, hookedDB),
		tickettypeCollection: busdb.NewCollection(tickettypeCollection, hookedDB),
		jobCollection:        busdb.NewCollection(jobCollection, hookedDB),
	}, nil
}

func setupDB(ctx context.Context, client driver.Client, dbName string) (driver.Database, error) {
	databaseExists, err := client.DatabaseExists(ctx, dbName)
	if err != nil {
		return nil, err
	}

	var db driver.Database
	if !databaseExists {
		db, err = client.CreateDatabase(ctx, dbName, nil)
	} else {
		db, err = client.Database(ctx, dbName)
	}
	if err != nil {
		return nil, err
	}

	collectionExists, err := db.CollectionExists(ctx, migrations.MigrationCollection)
	if err != nil {
		return nil, err
	}

	if !collectionExists {
		if _, err := db.CreateCollection(ctx, migrations.MigrationCollection, &driver.CreateCollectionOptions{
			KeyOptions: &driver.CollectionKeyOptions{AllowUserKeys: true},
		}); err != nil {
			log.Println(err)
		}
	}

	return db, nil
}
