package automation

import (
	"encoding/json"
	"fmt"
	"log"

	"golang.org/x/net/context"

	"github.com/SecurityBrewery/catalyst/bus"
	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/generated/models"
)

func jobAutomation(ctx context.Context, apiurl, apikey string, catalystBus *bus.Bus, db *database.Database) error {
	return catalystBus.SubscribeJob(func(automationMsg *bus.JobMsg) {
		job, err := db.JobCreate(ctx, automationMsg.ID, &models.JobForm{
			Automation: automationMsg.Automation,
			Payload:    automationMsg.Message.Payload,
			Origin:     automationMsg.Origin,
		})
		if err != nil {
			log.Println(err)
			return
		}

		automation, err := db.AutomationGet(ctx, automationMsg.Automation)
		if err != nil {
			log.Println(err)
			return
		}

		if automation.Script == "" {
			log.Println("automation is empty")
			return
		}

		if automationMsg.Message.Secrets == nil {
			automationMsg.Message.Secrets = map[string]string{}
		}
		automationMsg.Message.Secrets["catalyst_apikey"] = apikey
		automationMsg.Message.Secrets["catalyst_apiurl"] = apiurl

		scriptMessage, _ := json.Marshal(automationMsg.Message)

		containerID, logs, err := createContainer(ctx, automation.Image, automation.Script, string(scriptMessage))
		if err != nil {
			log.Println(err)
			return
		}

		if _, err := db.JobUpdate(ctx, automationMsg.ID, &models.Job{
			Automation: job.Automation,
			Container:  &containerID,
			Origin:     job.Origin,
			Output:     job.Output,
			Log:        &logs,
			Payload:    job.Payload,
			Status:     job.Status,
		}); err != nil {
			log.Println(err)
			return
		}

		var result map[string]interface{}

		stdout, _, err := runDocker(ctx, automationMsg.ID, containerID, db)
		if err != nil {
			result = map[string]interface{}{"error": fmt.Sprintf("error running script %s %s", err, string(stdout))}
		} else {
			var data map[string]interface{}
			if err := json.Unmarshal(stdout, &data); err != nil {
				result = map[string]interface{}{"error": string(stdout)}
			} else {
				result = data
			}
		}

		if err := catalystBus.PublishResult(automationMsg.Automation, result, automationMsg.Origin); err != nil {
			log.Println(err)
		}

		if err := db.JobComplete(ctx, automationMsg.ID, result); err != nil {
			log.Println(err)
			return
		}
	})
}

/*
func getAutomation(automationID string, config *Config) (*models.AutomationResponse, error) {
	req, err := http.NewRequest(http.MethodGet, config.CatalystAPIUrl+"/automations/"+automationID, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", config.CatalystAPIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var automation models.AutomationResponse
	if err := json.Unmarshal(b, &automation); err != nil {
		return nil, err
	}
	return &automation, nil
}
*/
