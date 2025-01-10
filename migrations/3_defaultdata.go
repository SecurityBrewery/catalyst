package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
)

func defaultDataUp(app core.App) error {
	for _, records := range [][]*core.Record{typeRecords(app)} {
		for _, record := range records {
			if err := app.SaveNoValidate(record); err != nil {
				return err
			}
		}
	}

	return nil
}

func typeRecords(app core.App) []*core.Record {
	collection, err := app.FindCollectionByNameOrId(TypeCollectionName)
	if err != nil {
		panic(err)
	}

	var records []*core.Record

	record := core.NewRecord(collection)
	record.Id = "incident"
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

	record = core.NewRecord(collection)
	record.Id = "alert"
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
