package schedule

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/multierr"

	"github.com/SecurityBrewery/catalyst/app"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/reaction/action"
)

type Schedule struct {
	Expression string `json:"expression"`
}

func Start(pb *app.App) {
	s, err := gocron.NewScheduler()
	if err != nil {
		slog.Error(fmt.Sprintf("failed to create scheduler: %v", err))
	}

	// add a job to the scheduler
	_, err = s.NewJob(
		gocron.CronJob(
			"* * * * *",
			false,
		),
		gocron.NewTask(
			func(ctx context.Context) {
				if err := runSchedule(ctx, pb); err != nil {
					slog.ErrorContext(ctx, "failed to run schedule", "error", err.Error())
				}
			},
		),
	)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to add cron job: %v", err))
	}

	s.Start()
}

func runSchedule(ctx context.Context, app *app.App) error {
	var errs error

	records, err := findByScheduleTrigger(ctx, app.Queries)
	if err != nil {
		errs = multierr.Append(errs, fmt.Errorf("failed to find schedule reaction: %w", err))
	}

	if len(records) == 0 {
		return nil
	}

	for _, hook := range records {
		_, err = action.Run(ctx, app, hook.Action, hook.Actiondata, "{}")
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to run hook reaction: %w", err))
		}
	}

	return errs
}

func findByScheduleTrigger(ctx context.Context, queries *sqlc.Queries) ([]*sqlc.ListReactionsRow, error) {
	reactions, err := queries.ListReactions(ctx, sqlc.ListReactionsParams{
		Offset: 0,
		Limit:  100,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find schedule reaction: %w", err)
	}

	if len(reactions) == 0 {
		return nil, nil
	}

	var errs error

	var matchedRecords []*sqlc.ListReactionsRow

	for _, reaction := range reactions {
		if reaction.Trigger != "schedule" {
			continue
		}

		var schedule Schedule
		if err := json.Unmarshal([]byte(reaction.Triggerdata), &schedule); err != nil {
			errs = multierr.Append(errs, err)

			continue
		}

		// if s.IsDue(moment) { TODO
		matchedRecords = append(matchedRecords, &reaction)
	}

	return matchedRecords, errs
}
