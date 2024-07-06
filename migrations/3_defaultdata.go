package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/security"
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
	record.SetId("y_" + security.PseudorandomString(5))
	record.Set("singular", "Incident")
	record.Set("plural", "Incidents")
	record.Set("icon", "Flame")
	record.Set("schema", `{"type":"object","properties":{"tlp":{"title":"TLP","type":"string"}}}`)

	records = append(records, record)

	record = models.NewRecord(collection)
	record.SetId("y_" + security.PseudorandomString(5))
	record.Set("singular", "Alert")
	record.Set("plural", "Alerts")
	record.Set("icon", "AlertTriangle")
	record.Set("schema", `{"type":"object","properties":{"severity":{"title":"Severity","type":"string"}},"required": ["severity"]}`)

	records = append(records, record)

	return records
}
