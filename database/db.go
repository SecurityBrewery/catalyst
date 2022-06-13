package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"

	"github.com/SecurityBrewery/catalyst/bus"
	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/database/migrations"
	"github.com/SecurityBrewery/catalyst/generated/model"
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
	SettingsCollectionName   = "settings"
	DashboardCollectionName  = "dashboards"

	TicketArtifactsGraphName     = "Graph"
	RelatedTicketsCollectionName = "related"
)

type Database struct {
	*busdb.BusDatabase
	Index *index.Index
	bus   *bus.Bus
	Hooks *hooks.Hooks

	templateCollection   *busdb.Collection[model.TicketTemplate]
	ticketCollection     *busdb.Collection[model.Ticket]
	playbookCollection   *busdb.Collection[model.PlaybookTemplate]
	automationCollection *busdb.Collection[model.Automation]
	userdataCollection   *busdb.Collection[model.UserData]
	userCollection       *busdb.Collection[model.User]
	tickettypeCollection *busdb.Collection[model.TicketType]
	jobCollection        *busdb.Collection[model.Job]
	settingsCollection   *busdb.Collection[model.Settings]
	dashboardCollection  *busdb.Collection[model.Dashboard]

	relatedCollection *busdb.Collection[driver.EdgeDocument]
	// containsCollection *busdb.Collection
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

	var err error
	var client driver.Client
	for {
		deadline, ok := ctx.Deadline()
		if ok && time.Until(deadline) < 0 {
			return nil, context.DeadlineExceeded
		}

		client, err = getClient(ctx, config)
		if err == nil {
			break
		}

		if errors.Is(err, context.DeadlineExceeded) {
			return nil, errors.New("could not load database, connection timed out")
		}

		log.Printf("could not connect to database: %s, retrying in 10 seconds\n", err)
		time.Sleep(time.Second * 10)
	}

	hooks.DatabaseAfterConnect(ctx, client, name)

	arangoDB, err := SetupDB(ctx, client, name)
	if err != nil {
		return nil, fmt.Errorf("DB setup failed: %w", err)
	}

	if err = migrations.PerformMigrations(ctx, arangoDB); err != nil {
		return nil, fmt.Errorf("migrations failed: %w", err)
	}

	ticketCollection, err := arangoDB.Collection(ctx, TicketCollectionName)
	if err != nil {
		return nil, err
	}
	templateCollection, err := arangoDB.Collection(ctx, TemplateCollectionName)
	if err != nil {
		return nil, err
	}
	playbookCollection, err := arangoDB.Collection(ctx, PlaybookCollectionName)
	if err != nil {
		return nil, err
	}
	relatedCollection, err := arangoDB.Collection(ctx, RelatedTicketsCollectionName)
	if err != nil {
		return nil, err
	}
	automationCollection, err := arangoDB.Collection(ctx, AutomationCollectionName)
	if err != nil {
		return nil, err
	}
	userdataCollection, err := arangoDB.Collection(ctx, UserDataCollectionName)
	if err != nil {
		return nil, err
	}
	userCollection, err := arangoDB.Collection(ctx, UserCollectionName)
	if err != nil {
		return nil, err
	}
	tickettypeCollection, err := arangoDB.Collection(ctx, TicketTypeCollectionName)
	if err != nil {
		return nil, err
	}
	jobCollection, err := arangoDB.Collection(ctx, JobCollectionName)
	if err != nil {
		return nil, err
	}
	settingsCollection, err := arangoDB.Collection(ctx, SettingsCollectionName)
	if err != nil {
		return nil, err
	}
	dashboardCollection, err := arangoDB.Collection(ctx, DashboardCollectionName)
	if err != nil {
		return nil, err
	}

	hookedDB, err := busdb.NewDatabase(ctx, arangoDB, bus)
	if err != nil {
		return nil, err
	}

	db := &Database{
		BusDatabase:          hookedDB,
		bus:                  bus,
		Index:                index,
		Hooks:                hooks,
		templateCollection:   busdb.NewCollection[model.TicketTemplate](templateCollection, hookedDB),
		ticketCollection:     busdb.NewCollection[model.Ticket](ticketCollection, hookedDB),
		playbookCollection:   busdb.NewCollection[model.PlaybookTemplate](playbookCollection, hookedDB),
		automationCollection: busdb.NewCollection[model.Automation](automationCollection, hookedDB),
		userdataCollection:   busdb.NewCollection[model.UserData](userdataCollection, hookedDB),
		userCollection:       busdb.NewCollection[model.User](userCollection, hookedDB),
		tickettypeCollection: busdb.NewCollection[model.TicketType](tickettypeCollection, hookedDB),
		jobCollection:        busdb.NewCollection[model.Job](jobCollection, hookedDB),
		settingsCollection:   busdb.NewCollection[model.Settings](settingsCollection, hookedDB),
		dashboardCollection:  busdb.NewCollection[model.Dashboard](dashboardCollection, hookedDB),
		relatedCollection:    busdb.NewCollection[driver.EdgeDocument](relatedCollection, hookedDB),
	}

	return db, nil
}

func getClient(ctx context.Context, config *Config) (driver.Client, error) {
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

	if _, err := client.Version(ctx); err != nil {
		return nil, err
	}

	return client, nil
}

func SetupDB(ctx context.Context, client driver.Client, dbName string) (driver.Database, error) {
	databaseExists, err := client.DatabaseExists(ctx, dbName)
	if err != nil {
		return nil, fmt.Errorf("could not check if database exists: %w", err)
	}

	var db driver.Database
	if !databaseExists {
		db, err = client.CreateDatabase(ctx, dbName, nil)
	} else {
		db, err = client.Database(ctx, dbName)
	}
	if err != nil {
		return nil, fmt.Errorf("could not create database: %w", err)
	}

	collectionExists, err := db.CollectionExists(ctx, migrations.MigrationCollection)
	if err != nil {
		return nil, fmt.Errorf("could not check if collection exists: %w", err)
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

func (db *Database) Truncate(ctx context.Context) {
	_ = db.templateCollection.Truncate(ctx)
	_ = db.ticketCollection.Truncate(ctx)
	_ = db.playbookCollection.Truncate(ctx)
	_ = db.automationCollection.Truncate(ctx)
	_ = db.userdataCollection.Truncate(ctx)
	_ = db.userCollection.Truncate(ctx)
	_ = db.tickettypeCollection.Truncate(ctx)
	_ = db.jobCollection.Truncate(ctx)
	_ = db.relatedCollection.Truncate(ctx)
	_ = db.settingsCollection.Truncate(ctx)
	_ = db.dashboardCollection.Truncate(ctx)
	// db.containsCollection.Truncate(ctx)
}
