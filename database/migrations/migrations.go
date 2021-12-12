package migrations

import (
	"context"
	"fmt"

	"github.com/arangodb/go-driver"

	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/models"
	"github.com/SecurityBrewery/catalyst/pointer"
)

const MigrationCollection string = "migrations"

type Migration interface {
	MID() string
	Migrate(ctx context.Context, driver driver.Database) error
}

func generateMigrations() ([]Migration, error) {
	// content here should never change
	return []Migration{
		&createCollection{ID: "create-log-collection", Name: "logs", DataType: "log", Schema: `{"properties":{"created":{"format":"date-time","type":"string"},"creator":{"type":"string"},"message":{"type":"string"},"reference":{"type":"string"}},"required":["created","creator","message","reference"],"type":"object"}`},
		&createCollection{ID: "create-ticket-collection", Name: "tickets", DataType: "ticket", Schema: `{"properties":{"artifacts":{"items":{"properties":{"enrichments":{"additionalProperties":{"properties":{"created":{"format":"date-time","type":"string"},"data":{"example":{"hash":"b7a067a742c20d07a7456646de89bc2d408a1153"},"properties":{},"type":"object"},"name":{"example":"hash.sha1","type":"string"}},"required":["created","data","name"],"type":"object"},"type":"object"},"name":{"example":"2.2.2.2","type":"string"},"status":{"example":"Unknown","type":"string"},"type":{"type":"string"}},"required":["name"],"type":"object"},"type":"array"},"comments":{"items":{"properties":{"created":{"format":"date-time","type":"string"},"creator":{"type":"string"},"message":{"type":"string"}},"required":["created","creator","message"],"type":"object"},"type":"array"},"created":{"format":"date-time","type":"string"},"details":{"example":{"description":"my little incident"},"properties":{},"type":"object"},"files":{"items":{"properties":{"key":{"example":"myfile","type":"string"},"name":{"example":"notes.docx","type":"string"}},"required":["key","name"],"type":"object"},"type":"array"},"modified":{"format":"date-time","type":"string"},"name":{"example":"WannyCry","type":"string"},"owner":{"example":"bob","type":"string"},"playbooks":{"additionalProperties":{"properties":{"name":{"example":"Phishing","type":"string"},"tasks":{"additionalProperties":{"properties":{"automation":{"type":"string"},"closed":{"format":"date-time","type":"string"},"created":{"format":"date-time","type":"string"},"data":{"properties":{},"type":"object"},"done":{"type":"boolean"},"join":{"example":false,"type":"boolean"},"payload":{"additionalProperties":{"type":"string"},"type":"object"},"name":{"example":"Inform user","type":"string"},"next":{"additionalProperties":{"type":"string"},"type":"object"},"owner":{"type":"string"},"schema":{"properties":{},"type":"object"},"type":{"enum":["task","input","automation"],"example":"task","type":"string"}},"required":["created","done","name","type"],"type":"object"},"type":"object"}},"required":["name","tasks"],"type":"object"},"type":"object"},"read":{"example":["bob"],"items":{"type":"string"},"type":"array"},"references":{"items":{"properties":{"href":{"example":"https://cve.mitre.org/cgi-bin/cvename.cgi?name=cve-2017-0144","type":"string"},"name":{"example":"CVE-2017-0144","type":"string"}},"required":["href","name"],"type":"object"},"type":"array"},"schema":{"example":"{}","type":"string"},"status":{"example":"open","type":"string"},"type":{"example":"incident","type":"string"},"write":{"example":["alice"],"items":{"type":"string"},"type":"array"}},"required":["created","modified","name","schema","status","type"],"type":"object"}`},
		&createCollection{ID: "create-template-collection", Name: "templates", DataType: "template", Schema: `{"properties":{"name":{"type":"string"},"schema":{"type":"string"}},"required":["name","schema"],"type":"object"}`},
		&createCollection{ID: "create-playbook-collection", Name: "playbooks", DataType: "playbook", Schema: `{"properties":{"name":{"type":"string"},"yaml":{"type":"string"}},"required":["name","yaml"],"type":"object"}`},
		&createCollection{ID: "create-automation-collection", Name: "automations", DataType: "automation", Schema: `{"properties":{"image":{"type":"string"},"script":{"type":"string"}},"required":["image","script"],"type":"object"}`},
		&createCollection{ID: "create-userdata-collection", Name: "userdata", DataType: "userdata", Schema: `{"properties":{"email":{"type":"string"},"image":{"type":"string"},"name":{"type":"string"},"timeformat":{"title":"Time Format (https://moment.github.io/luxon/docs/manual/formatting.html#table-of-tokens)","type":"string"}},"type":"object"}`},
		&createCollection{ID: "create-tickettype-collection", Name: "tickettypes", DataType: "tickettype", Schema: `{"properties":{"default_groups":{"items":{"type":"string"},"type":"array"},"default_playbooks":{"items":{"type":"string"},"type":"array"},"default_template":{"type":"string"},"icon":{"type":"string"},"name":{"type":"string"}},"required":["default_playbooks","default_template","icon","name"],"type":"object"}`},
		&createCollection{ID: "create-user-collection", Name: "users", DataType: "user", Schema: `{"properties":{"apikey":{"type":"boolean"},"blocked":{"type":"boolean"},"roles":{"items":{"type":"string"},"type":"array"},"sha256":{"type":"string"}},"required":["apikey","blocked","roles"],"type":"object"}`},

		&createGraph{ID: "create-ticket-graph", Name: "Graph", EdgeDefinitions: []driver.EdgeDefinition{{Collection: "related", From: []string{"tickets"}, To: []string{"tickets"}}}},

		&createDocument{ID: "create-template-default", Collection: "templates", Document: &busdb.Keyed{Key: "default", Doc: models.TicketTemplate{Schema: DefaultTemplateSchema, Name: "Default"}}},
		&createDocument{ID: "create-automation-vt.hash", Collection: "automations", Document: &busdb.Keyed{Key: "vt.hash", Doc: models.Automation{Image: "docker.io/python:3", Script: VTHashAutomation}}},
		&createDocument{ID: "create-automation-comment", Collection: "automations", Document: &busdb.Keyed{Key: "comment", Doc: models.Automation{Image: "docker.io/python:3", Script: CommentAutomation}}},
		&createDocument{ID: "create-automation-thehive", Collection: "automations", Document: &busdb.Keyed{Key: "thehive", Doc: models.Automation{Image: "docker.io/python:3", Script: TheHiveAutomation}}},
		&createDocument{ID: "create-automation-hash.sha1", Collection: "automations", Document: &busdb.Keyed{Key: "hash.sha1", Doc: models.Automation{Image: "docker.io/python:3", Script: SHA1HashAutomation}}},
		&createDocument{ID: "create-playbook-malware", Collection: "playbooks", Document: &busdb.Keyed{Key: "malware", Doc: models.PlaybookTemplate{Name: "Malware", Yaml: MalwarePlaybook}}},
		&createDocument{ID: "create-playbook-phishing", Collection: "playbooks", Document: &busdb.Keyed{Key: "phishing", Doc: models.PlaybookTemplate{Name: "Phishing", Yaml: PhishingPlaybook}}},
		&createDocument{ID: "create-tickettype-alert", Collection: "tickettypes", Document: &busdb.Keyed{Key: "alert", Doc: models.TicketType{Name: "Alerts", Icon: "mdi-alert", DefaultTemplate: "default", DefaultPlaybooks: []string{}, DefaultGroups: nil}}},
		&createDocument{ID: "create-tickettype-incident", Collection: "tickettypes", Document: &busdb.Keyed{Key: "incident", Doc: models.TicketType{Name: "Incidents", Icon: "mdi-radioactive", DefaultTemplate: "default", DefaultPlaybooks: []string{}, DefaultGroups: nil}}},
		&createDocument{ID: "create-tickettype-investigation", Collection: "tickettypes", Document: &busdb.Keyed{Key: "investigation", Doc: models.TicketType{Name: "Forensic Investigations", Icon: "mdi-fingerprint", DefaultTemplate: "default", DefaultPlaybooks: []string{}, DefaultGroups: nil}}},
		&createDocument{ID: "create-tickettype-hunt", Collection: "tickettypes", Document: &busdb.Keyed{Key: "hunt", Doc: models.TicketType{Name: "Threat Hunting", Icon: "mdi-target", DefaultTemplate: "default", DefaultPlaybooks: []string{}, DefaultGroups: nil}}},

		&updateSchema{ID: "update-automation-collection-1", Name: "automations", DataType: "automation", Schema: `{"properties":{"image":{"type":"string"},"script":{"type":"string"}},"required":["image","script"],"type":"object"}`},
		&updateDocument{ID: "update-automation-vt.hash-1", Collection: "automations", Key: "vt.hash", Document: models.Automation{Image: "docker.io/python:3", Script: VTHashAutomation, Schema: pointer.String(`{"title":"Input","type":"object","properties":{"default":{"type":"string","title":"Value"}},"required":["default"]}`), Type: []string{"global", "artifact", "playbook"}}},
		&updateDocument{ID: "update-automation-comment-1", Collection: "automations", Key: "comment", Document: models.Automation{Image: "docker.io/python:3", Script: CommentAutomation, Type: []string{"playbook"}}},
		&updateDocument{ID: "update-automation-thehive-1", Collection: "automations", Key: "thehive", Document: models.Automation{Image: "docker.io/python:3", Script: TheHiveAutomation, Schema: pointer.String(`{"title":"TheHive credentials","type":"object","properties":{"thehiveurl":{"type":"string","title":"TheHive URL (e.g. 'https://thehive.example.org')"},"thehivekey":{"type":"string","title":"TheHive API Key"},"skip_files":{"type":"boolean", "default": true, "title":"Skip Files (much faster)"},"keep_ids":{"type":"boolean", "default": true, "title":"Keep IDs and overwrite existing IDs"}},"required":["thehiveurl", "thehivekey", "skip_files", "keep_ids"]}`), Type: []string{"global"}}},
		&updateDocument{ID: "update-automation-hash.sha1-1", Collection: "automations", Key: "hash.sha1", Document: models.Automation{Image: "docker.io/python:3", Script: SHA1HashAutomation, Schema: pointer.String(`{"title":"Input","type":"object","properties":{"default":{"type":"string","title":"Value"}},"required":["default"]}`), Type: []string{"global", "artifact", "playbook"}}},

		&createCollection{ID: "create-job-collection", Name: "jobs", DataType: "job", Schema: `{"properties":{"automation":{"type":"string"},"log":{"type":"string"},"payload":{},"origin":{"properties":{"artifact_origin":{"properties":{"artifact":{"type":"string"},"ticket_id":{"format":"int64","type":"integer"}},"required":["artifact","ticket_id"],"type":"object"},"task_origin":{"properties":{"playbook_id":{"type":"string"},"task_id":{"type":"string"},"ticket_id":{"format":"int64","type":"integer"}},"required":["playbook_id","task_id","ticket_id"],"type":"object"}},"type":"object"},"output":{"properties":{},"type":"object"},"running":{"type":"boolean"},"status":{"type":"string"}},"required":["automation","running","status"],"type":"object"}`},
	}, nil
}

func loadSchema(dataType, jsonschema string) (*driver.CollectionSchemaOptions, error) {
	ticketCollectionSchema := &driver.CollectionSchemaOptions{Level: driver.CollectionSchemaLevelStrict, Message: fmt.Sprintf("Validation of %s failed", dataType)}

	err := ticketCollectionSchema.LoadRule([]byte(jsonschema))
	return ticketCollectionSchema, err
}

type migration struct {
	Key string `json:"_key"`
}

func PerformMigrations(ctx context.Context, db driver.Database) error {
	collection, err := db.Collection(ctx, MigrationCollection)
	if err != nil {
		return err
	}

	migrations, err := generateMigrations()
	if err != nil {
		return fmt.Errorf("could not generate migrations: %w", err)
	}

	for _, m := range migrations {
		migrationRan, err := collection.DocumentExists(ctx, m.MID())
		if err != nil {
			return err
		}

		if !migrationRan {
			if err := m.Migrate(ctx, db); err != nil {
				return fmt.Errorf("migration %s failed: %w", m.MID(), err)
			}

			if _, err := collection.CreateDocument(ctx, &migration{Key: m.MID()}); err != nil {
				return fmt.Errorf("could not save %s migration document: %w", m.MID(), err)
			}
		}
	}
	return nil
}

type createCollection struct {
	ID       string
	Name     string
	DataType string
	Schema   string
}

func (m *createCollection) MID() string {
	return m.ID
}

func (m *createCollection) Migrate(ctx context.Context, db driver.Database) error {
	schema, err := loadSchema(m.DataType, m.Schema)
	if err != nil {
		return err
	}

	_, err = db.CreateCollection(ctx, m.Name, &driver.CreateCollectionOptions{
		Schema: schema,
	})

	return err
}

type updateSchema struct {
	ID       string
	Name     string
	DataType string
	Schema   string
}

func (m *updateSchema) MID() string {
	return m.ID
}

func (m *updateSchema) Migrate(ctx context.Context, db driver.Database) error {
	schema, err := loadSchema(m.DataType, m.Schema)
	if err != nil {
		return err
	}

	col, err := db.Collection(ctx, m.Name)
	if err != nil {
		return err
	}

	err = col.SetProperties(ctx, driver.SetCollectionPropertiesOptions{
		Schema: schema,
	})

	return err
}

type createGraph struct {
	ID              string
	Name            string
	EdgeDefinitions []driver.EdgeDefinition
}

func (m *createGraph) MID() string {
	return m.ID
}

func (m *createGraph) Migrate(ctx context.Context, db driver.Database) error {
	_, err := db.CreateGraph(ctx, m.Name, &driver.CreateGraphOptions{
		EdgeDefinitions: m.EdgeDefinitions,
	})
	return err
}

type createDocument struct {
	ID         string
	Collection string
	Document   interface{}
}

func (m *createDocument) MID() string {
	return m.ID
}

func (m *createDocument) Migrate(ctx context.Context, driver driver.Database) error {
	collection, err := driver.Collection(ctx, m.Collection)
	if err != nil {
		return err
	}

	_, err = collection.CreateDocument(ctx, m.Document)
	return err
}

type updateDocument struct {
	ID         string
	Collection string
	Key        string
	Document   interface{}
}

func (m *updateDocument) MID() string {
	return m.ID
}

func (m *updateDocument) Migrate(ctx context.Context, driver driver.Database) error {
	collection, err := driver.Collection(ctx, m.Collection)
	if err != nil {
		return err
	}

	exists, err := collection.DocumentExists(ctx, m.Key)
	if err != nil {
		return err
	}

	if !exists {
		_, err = collection.CreateDocument(ctx, m.Document)
		return err
	}

	_, err = collection.ReplaceDocument(ctx, m.Key, m.Document)
	return err
}
