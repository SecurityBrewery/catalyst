package schedule

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/cron"
	"go.uber.org/multierr"

	"github.com/SecurityBrewery/catalyst/migrations"
	"github.com/SecurityBrewery/catalyst/reaction/action"
)

type Schedule struct {
	Expression string `json:"expression"`
}

func Start(pb *pocketbase.PocketBase) {
	scheduler := cron.New()

	if err := scheduler.Add("reactions", "* * * * *", func() {
		ctx := context.Background()

		moment := cron.NewMoment(time.Now())

		if err := runSchedule(ctx, pb.App, moment); err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("failed to run hook reaction: %v", err))
		}
	}); err != nil {
		slog.Error(fmt.Sprintf("failed to add cron job: %v", err))
	}

	scheduler.Start()
}

func runSchedule(ctx context.Context, app core.App, moment *cron.Moment) error {
	var errs error

	records, err := findByScheduleTrigger(app.Dao(), moment)
	if err != nil {
		errs = multierr.Append(errs, fmt.Errorf("failed to find schedule reaction: %w", err))
	}

	if len(records) == 0 {
		return nil
	}

	for _, hook := range records {
		_, err = action.Run(ctx, app, hook.GetString("action"), hook.GetString("actiondata"), "{}")
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to run hook reaction: %w", err))
		}
	}

	return errs
}

func findByScheduleTrigger(dao *daos.Dao, moment *cron.Moment) ([]*models.Record, error) {
	records, err := dao.FindRecordsByExpr(migrations.ReactionCollectionName, dbx.HashExp{"trigger": "schedule"})
	if err != nil {
		return nil, fmt.Errorf("failed to find schedule reaction: %w", err)
	}

	if len(records) == 0 {
		return nil, nil
	}

	var errs error

	var matchedRecords []*models.Record

	for _, record := range records {
		var schedule Schedule
		if err := json.Unmarshal([]byte(record.GetString("triggerdata")), &schedule); err != nil {
			errs = multierr.Append(errs, err)

			continue
		}

		s, err := cron.NewSchedule(schedule.Expression)
		if err != nil {
			errs = multierr.Append(errs, err)

			continue
		}

		if s.IsDue(moment) {
			matchedRecords = append(matchedRecords, record)
		}
	}

	return matchedRecords, errs
}
