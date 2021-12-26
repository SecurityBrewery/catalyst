package database

import (
	"context"
	"errors"

	"github.com/arangodb/go-driver"
	"github.com/iancoleman/strcase"
	"github.com/icza/dyno"
	"gopkg.in/yaml.v3"

	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/models"
	"github.com/SecurityBrewery/catalyst/time"
)

type PlaybookYAML struct {
	Name  string              `yaml:"name"`
	Tasks map[string]TaskYAML `yaml:"tasks"`
}

type TaskYAML struct {
	Name       string            `yaml:"name"`
	Type       string            `yaml:"type"`
	Schema     interface{}       `yaml:"schema"`
	Automation string            `yaml:"automation"`
	Payload    map[string]string `yaml:"payload"`
	Next       map[string]string `yaml:"next"`
	Join       bool              `yaml:"join"`
}

func toPlaybooks(docs []*models.PlaybookTemplateForm) (map[string]*models.Playbook, error) {
	playbooks := map[string]*models.Playbook{}
	for _, doc := range docs {
		playbook, err := toPlaybook(doc)
		if err != nil {
			return nil, err
		}
		if doc.ID != nil {
			playbooks[*doc.ID] = playbook
		} else {
			playbooks[strcase.ToKebab(playbook.Name)] = playbook
		}
	}
	return playbooks, nil
}

func toPlaybook(doc *models.PlaybookTemplateForm) (*models.Playbook, error) {
	ticketPlaybook := &models.Playbook{}
	err := yaml.Unmarshal([]byte(doc.Yaml), ticketPlaybook)
	if err != nil {
		return nil, err
	}
	for idx, task := range ticketPlaybook.Tasks {
		if task.Schema != nil {
			task.Schema = dyno.ConvertMapI2MapS(task.Schema.(map[string]interface{}))
		}
		task.Created = time.Now().UTC()
		ticketPlaybook.Tasks[idx] = task
	}
	return ticketPlaybook, nil
}

func toPlaybookTemplateResponse(key string, doc *models.PlaybookTemplate) *models.PlaybookTemplateResponse {
	return &models.PlaybookTemplateResponse{ID: key, Name: doc.Name, Yaml: doc.Yaml}
}

func (db *Database) PlaybookCreate(ctx context.Context, playbook *models.PlaybookTemplateForm) (*models.PlaybookTemplateResponse, error) {
	if playbook == nil {
		return nil, errors.New("requires playbook")
	}

	var playbookYAML PlaybookYAML
	err := yaml.Unmarshal([]byte(playbook.Yaml), &playbookYAML)
	if err != nil {
		return nil, err
	}

	if playbookYAML.Name == "" {
		return nil, errors.New("requires template name")
	}
	p := models.PlaybookTemplate{Name: playbookYAML.Name, Yaml: playbook.Yaml}

	var doc models.PlaybookTemplate
	newctx := driver.WithReturnNew(ctx, &doc)

	meta, err := db.playbookCollection.CreateDocument(ctx, newctx, strcase.ToKebab(playbookYAML.Name), p)
	if err != nil {
		return nil, err
	}

	return toPlaybookTemplateResponse(meta.Key, &doc), nil
}

func (db *Database) PlaybookGet(ctx context.Context, id string) (*models.PlaybookTemplateResponse, error) {
	doc := models.PlaybookTemplate{}
	meta, err := db.playbookCollection.ReadDocument(ctx, id, &doc)
	if err != nil {
		return nil, err
	}

	return toPlaybookTemplateResponse(meta.Key, &doc), nil
}

func (db *Database) PlaybookDelete(ctx context.Context, id string) error {
	_, err := db.playbookCollection.RemoveDocument(ctx, id)
	return err
}

func (db *Database) PlaybookUpdate(ctx context.Context, id string, playbook *models.PlaybookTemplateForm) (*models.PlaybookTemplateResponse, error) {
	var pb PlaybookYAML
	err := yaml.Unmarshal([]byte(playbook.Yaml), &pb)
	if err != nil {
		return nil, err
	}

	if pb.Name == "" {
		return nil, errors.New("requires template name")
	}

	var doc models.PlaybookTemplate
	ctx = driver.WithReturnNew(ctx, &doc)

	meta, err := db.playbookCollection.ReplaceDocument(ctx, id, models.PlaybookTemplate{Name: pb.Name, Yaml: playbook.Yaml})
	if err != nil {
		return nil, err
	}

	return toPlaybookTemplateResponse(meta.Key, &doc), nil
}

func (db *Database) PlaybookList(ctx context.Context) ([]*models.PlaybookTemplateResponse, error) {
	query := "FOR d IN @@collection RETURN d"
	cursor, _, err := db.Query(ctx, query, map[string]interface{}{"@collection": PlaybookCollectionName}, busdb.ReadOperation)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	var docs []*models.PlaybookTemplateResponse
	for {
		var doc models.PlaybookTemplate
		meta, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}
		docs = append(docs, toPlaybookTemplateResponse(meta.Key, &doc))
	}

	return docs, err
}
