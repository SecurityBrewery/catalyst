package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

const (
	AttackIntegrationID        = "attack"
	VulnerabilityIntegrationID = "vulnerability"
)

func integrationsDataUp(db dbx.Builder) error {
	dao := daos.New(db)

	collection, err := dao.FindCollectionByNameOrId(IntegrationCollectionName)
	if err != nil {
		return err
	}

	record := models.NewRecord(collection)
	record.SetId(AttackIntegrationID)
	record.Set("name", "attack")
	record.Set("plugin", "attack")

	if err := dao.SaveRecord(record); err != nil {
		return err
	}

	record = models.NewRecord(collection)
	record.SetId(VulnerabilityIntegrationID)
	record.Set("name", "vulnerability")
	record.Set("plugin", "vulnerability")

	return dao.SaveRecord(record)
}

func integrationsDataDown(db dbx.Builder) error {
	dao := daos.New(db)

	record, err := dao.FindRecordById(IntegrationCollectionName, AttackIntegrationID)
	if err != nil {
		return err
	}

	if err := dao.DeleteRecord(record); err != nil {
		return err
	}

	record, err = dao.FindRecordById(IntegrationCollectionName, VulnerabilityIntegrationID)
	if err != nil {
		return err
	}

	return dao.DeleteRecord(record)
}
