package migrations

import (
	"context"
	"fmt"

	"github.com/arangodb/go-driver"

	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/model"
	"github.com/SecurityBrewery/catalyst/generated/pointer"
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

		&createDocument[busdb.Keyed[model.TicketTemplate]]{ID: "create-template-default", Collection: "templates", Document: &busdb.Keyed[model.TicketTemplate]{Key: "default", Doc: &model.TicketTemplate{Schema: DefaultTemplateSchema, Name: "Default"}}},
		&createDocument[busdb.Keyed[model.Automation]]{ID: "create-automation-vt.hash", Collection: "automations", Document: &busdb.Keyed[model.Automation]{Key: "vt.hash", Doc: &model.Automation{Image: "docker.io/python:3", Script: VTHashAutomation}}},
		&createDocument[busdb.Keyed[model.Automation]]{ID: "create-automation-comment", Collection: "automations", Document: &busdb.Keyed[model.Automation]{Key: "comment", Doc: &model.Automation{Image: "docker.io/python:3", Script: CommentAutomation}}},
		&createDocument[busdb.Keyed[model.Automation]]{ID: "create-automation-hash.sha1", Collection: "automations", Document: &busdb.Keyed[model.Automation]{Key: "hash.sha1", Doc: &model.Automation{Image: "docker.io/python:3", Script: SHA1HashAutomation}}},
		&createDocument[busdb.Keyed[model.PlaybookTemplate]]{ID: "create-playbook-malware", Collection: "playbooks", Document: &busdb.Keyed[model.PlaybookTemplate]{Key: "malware", Doc: &model.PlaybookTemplate{Name: "Malware", Yaml: MalwarePlaybook}}},
		&createDocument[busdb.Keyed[model.PlaybookTemplate]]{ID: "create-playbook-phishing", Collection: "playbooks", Document: &busdb.Keyed[model.PlaybookTemplate]{Key: "phishing", Doc: &model.PlaybookTemplate{Name: "Phishing", Yaml: PhishingPlaybook}}},
		&createDocument[busdb.Keyed[model.TicketType]]{ID: "create-tickettype-alert", Collection: "tickettypes", Document: &busdb.Keyed[model.TicketType]{Key: "alert", Doc: &model.TicketType{Name: "Alerts", Icon: "mdi-alert", DefaultTemplate: "default", DefaultPlaybooks: []string{}, DefaultGroups: nil}}},
		&createDocument[busdb.Keyed[model.TicketType]]{ID: "create-tickettype-incident", Collection: "tickettypes", Document: &busdb.Keyed[model.TicketType]{Key: "incident", Doc: &model.TicketType{Name: "Incidents", Icon: "mdi-radioactive", DefaultTemplate: "default", DefaultPlaybooks: []string{}, DefaultGroups: nil}}},
		&createDocument[busdb.Keyed[model.TicketType]]{ID: "create-tickettype-investigation", Collection: "tickettypes", Document: &busdb.Keyed[model.TicketType]{Key: "investigation", Doc: &model.TicketType{Name: "Forensic Investigations", Icon: "mdi-fingerprint", DefaultTemplate: "default", DefaultPlaybooks: []string{}, DefaultGroups: nil}}},
		&createDocument[busdb.Keyed[model.TicketType]]{ID: "create-tickettype-hunt", Collection: "tickettypes", Document: &busdb.Keyed[model.TicketType]{Key: "hunt", Doc: &model.TicketType{Name: "Threat Hunting", Icon: "mdi-target", DefaultTemplate: "default", DefaultPlaybooks: []string{}, DefaultGroups: nil}}},

		&updateSchema{ID: "update-automation-collection-1", Name: "automations", DataType: "automation", Schema: `{"properties":{"image":{"type":"string"},"script":{"type":"string"}},"required":["image","script"],"type":"object"}`},
		&updateDocument[model.Automation]{ID: "update-automation-vt.hash-1", Collection: "automations", Key: "vt.hash", Document: &model.Automation{Image: "docker.io/python:3", Script: VTHashAutomation, Schema: pointer.String(`{"title":"Input","type":"object","properties":{"default":{"type":"string","title":"Value"}},"required":["default"]}`), Type: []string{"global", "artifact", "playbook"}}},
		&updateDocument[model.Automation]{ID: "update-automation-comment-1", Collection: "automations", Key: "comment", Document: &model.Automation{Image: "docker.io/python:3", Script: CommentAutomation, Type: []string{"playbook"}}},
		&updateDocument[model.Automation]{ID: "update-automation-hash.sha1-1", Collection: "automations", Key: "hash.sha1", Document: &model.Automation{Image: "docker.io/python:3", Script: SHA1HashAutomation, Schema: pointer.String(`{"title":"Input","type":"object","properties":{"default":{"type":"string","title":"Value"}},"required":["default"]}`), Type: []string{"global", "artifact", "playbook"}}},

		&createCollection{ID: "create-job-collection", Name: "jobs", DataType: "job", Schema: `{"properties":{"automation":{"type":"string"},"log":{"type":"string"},"payload":{},"origin":{"properties":{"artifact_origin":{"properties":{"artifact":{"type":"string"},"ticket_id":{"format":"int64","type":"integer"}},"required":["artifact","ticket_id"],"type":"object"},"task_origin":{"properties":{"playbook_id":{"type":"string"},"task_id":{"type":"string"},"ticket_id":{"format":"int64","type":"integer"}},"required":["playbook_id","task_id","ticket_id"],"type":"object"}},"type":"object"},"output":{"properties":{},"type":"object"},"running":{"type":"boolean"},"status":{"type":"string"}},"required":["automation","running","status"],"type":"object"}`},

		&createDocument[busdb.Keyed[model.PlaybookTemplate]]{ID: "create-playbook-simple", Collection: "playbooks", Document: &busdb.Keyed[model.PlaybookTemplate]{Key: "simple", Doc: &model.PlaybookTemplate{Name: "Simple", Yaml: SimplePlaybook}}},

		&createCollection{ID: "create-settings-collection", Name: "settings", DataType: "settings", Schema: `{"type":"object","properties":{"artifactStates":{"title":"Artifact States","items":{"type":"object","properties":{"color":{"title":"Color","type":"string","enum":["error","info","success","warning"]},"icon":{"title":"Icon (https://materialdesignicons.com)","type":"string"},"id":{"title":"ID","type":"string"},"name":{"title":"Name","type":"string"}},"required":["id","name","icon"]},"type":"array"},"artifactKinds":{"title":"Artifact Kinds","items":{"type":"object","properties":{"color":{"title":"Color","type":"string","enum":["error","info","success","warning"]},"icon":{"title":"Icon (https://materialdesignicons.com)","type":"string"},"id":{"title":"ID","type":"string"},"name":{"title":"Name","type":"string"}},"required":["id","name","icon"]},"type":"array"},"timeformat":{"title":"Time Format","type":"string"}},"required":["timeformat","artifactKinds","artifactStates"]}`},
		&createDocument[busdb.Keyed[model.Settings]]{ID: "create-settings-global", Collection: "settings", Document: &busdb.Keyed[model.Settings]{Key: "global", Doc: &model.Settings{ArtifactStates: []*model.Type{{Icon: "mdi-help-circle-outline", ID: "unknown", Name: "Unknown", Color: pointer.String(model.TypeColorInfo)}, {Icon: "mdi-skull", ID: "malicious", Name: "Malicious", Color: pointer.String(model.TypeColorError)}, {Icon: "mdi-check", ID: "clean", Name: "Clean", Color: pointer.String(model.TypeColorSuccess)}}, ArtifactKinds: []*model.Type{{Icon: "mdi-server", ID: "asset", Name: "Asset"}, {Icon: "mdi-bullseye", ID: "ioc", Name: "IOC"}}, Timeformat: "YYYY-MM-DDThh:mm:ss"}}},

		&updateSchema{ID: "update-ticket-collection", Name: "tickets", DataType: "ticket", Schema: `{"properties":{"artifacts":{"items":{"properties":{"enrichments":{"additionalProperties":{"properties":{"created":{"format":"date-time","type":"string"},"data":{"example":{"hash":"b7a067a742c20d07a7456646de89bc2d408a1153"},"properties":{},"type":"object"},"name":{"example":"hash.sha1","type":"string"}},"required":["created","data","name"],"type":"object"},"type":"object"},"name":{"example":"2.2.2.2","type":"string"},"status":{"example":"Unknown","type":"string"},"type":{"type":"string"},"kind":{"type":"string"}},"required":["name"],"type":"object"},"type":"array"},"comments":{"items":{"properties":{"created":{"format":"date-time","type":"string"},"creator":{"type":"string"},"message":{"type":"string"}},"required":["created","creator","message"],"type":"object"},"type":"array"},"created":{"format":"date-time","type":"string"},"details":{"example":{"description":"my little incident"},"properties":{},"type":"object"},"files":{"items":{"properties":{"key":{"example":"myfile","type":"string"},"name":{"example":"notes.docx","type":"string"}},"required":["key","name"],"type":"object"},"type":"array"},"modified":{"format":"date-time","type":"string"},"name":{"example":"WannyCry","type":"string"},"owner":{"example":"bob","type":"string"},"playbooks":{"additionalProperties":{"properties":{"name":{"example":"Phishing","type":"string"},"tasks":{"additionalProperties":{"properties":{"automation":{"type":"string"},"closed":{"format":"date-time","type":"string"},"created":{"format":"date-time","type":"string"},"data":{"properties":{},"type":"object"},"done":{"type":"boolean"},"join":{"example":false,"type":"boolean"},"payload":{"additionalProperties":{"type":"string"},"type":"object"},"name":{"example":"Inform user","type":"string"},"next":{"additionalProperties":{"type":"string"},"type":"object"},"owner":{"type":"string"},"schema":{"properties":{},"type":"object"},"type":{"enum":["task","input","automation"],"example":"task","type":"string"}},"required":["created","done","name","type"],"type":"object"},"type":"object"}},"required":["name","tasks"],"type":"object"},"type":"object"},"read":{"example":["bob"],"items":{"type":"string"},"type":"array"},"references":{"items":{"properties":{"href":{"example":"https://cve.mitre.org/cgi-bin/cvename.cgi?name=cve-2017-0144","type":"string"},"name":{"example":"CVE-2017-0144","type":"string"}},"required":["href","name"],"type":"object"},"type":"array"},"schema":{"example":"{}","type":"string"},"status":{"example":"open","type":"string"},"type":{"example":"incident","type":"string"},"write":{"example":["alice"],"items":{"type":"string"},"type":"array"}},"required":["created","modified","name","schema","status","type"],"type":"object"}`},

		&createCollection{ID: "create-dashboard-collection", Name: "dashboards", DataType: "dashboards", Schema: `{"type":"object","properties":{"name":{"type":"string"},"widgets":{"items":{"type":"object","properties":{"aggregation":{"type":"string"},"filter":{"type":"string"},"name":{"type":"string"},"type":{"enum":[ "bar", "line", "pie" ]},"width": { "type": "integer", "minimum": 1, "maximum": 12 }},"required":["name","aggregation", "type", "width"]},"type":"array"}},"required":["name","widgets"]}`},

		&updateDocument[model.Settings]{ID: "update-settings-global-1", Collection: "settings", Key: "global", Document: &model.Settings{ArtifactStates: []*model.Type{{Icon: "mdi-help-circle-outline", ID: "unknown", Name: "Unknown", Color: pointer.String(model.TypeColorInfo)}, {Icon: "mdi-skull", ID: "malicious", Name: "Malicious", Color: pointer.String(model.TypeColorError)}, {Icon: "mdi-check", ID: "clean", Name: "Clean", Color: pointer.String(model.TypeColorSuccess)}}, ArtifactKinds: []*model.Type{{Icon: "mdi-server", ID: "asset", Name: "Asset"}, {Icon: "mdi-bullseye", ID: "ioc", Name: "IOC"}}, Timeformat: "yyyy-MM-dd hh:mm:ss"}},
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

type createDocument[T any] struct {
	ID         string
	Collection string
	Document   *T
}

func (m *createDocument[T]) MID() string {
	return m.ID
}

func (m *createDocument[T]) Migrate(ctx context.Context, driver driver.Database) error {
	collection, err := driver.Collection(ctx, m.Collection)
	if err != nil {
		return err
	}

	_, err = collection.CreateDocument(ctx, m.Document)

	return err
}

type updateDocument[T any] struct {
	ID         string
	Collection string
	Key        string
	Document   *T
}

func (m *updateDocument[T]) MID() string {
	return m.ID
}

func (m *updateDocument[T]) Migrate(ctx context.Context, driver driver.Database) error {
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
