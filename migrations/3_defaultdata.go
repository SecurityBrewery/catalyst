package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

func defaultDataUp(db dbx.Builder) error {
	dao := daos.New(db)

	for _, records := range [][]*models.Record{typeRecords(dao)} {
		for _, record := range records {
			if err := dao.SaveRecord(record); err != nil {
				return err
			}
		}
	}

	return nil
}

func typeRecords(dao *daos.Dao) []*models.Record {
	collection, err := dao.FindCollectionByNameOrId(TypeCollectionName)
	if err != nil {
		panic(err)
	}

	var records []*models.Record

	record := models.NewRecord(collection)
	record.SetId("incident")
	record.Set("singular", "Incident")
	record.Set("plural", "Incidents")
	record.Set("icon", "Flame")
	record.Set("schema", s(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"severity": map[string]any{
				"title": "Severity",
				"enum":  []string{"Low", "Medium", "High"},
			},
		},
		"required": []string{"severity"},
	}))

	records = append(records, record)

	record = models.NewRecord(collection)
	record.SetId("alert")
	record.Set("singular", "Alert")
	record.Set("plural", "Alerts")
	record.Set("icon", "AlertTriangle")
	record.Set("schema", s(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"severity": map[string]any{
				"title": "Severity",
				"enum":  []string{"Low", "Medium", "High"},
			},
		},
		"required": []string{"severity"},
	}))

	records = append(records, record)

	return records
}

func s(m map[string]any) string {
	b, _ := json.Marshal(m) //nolint:errchkjson

	return string(b)
}
