package schedule

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-co-op/gocron/v2"

	"github.com/SecurityBrewery/catalyst/app/auth"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
	"github.com/SecurityBrewery/catalyst/app/reaction/action"
)

type Scheduler struct {
	scheduler gocron.Scheduler
	queries   *sqlc.Queries
	config    *auth.Config
	auth      *auth.Service
}

type Schedule struct {
	Expression string `json:"expression"`
}

func New(ctx context.Context, config *auth.Config, service *auth.Service, queries *sqlc.Queries) (*Scheduler, error) {
	innerScheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	scheduler := &Scheduler{
		config:    config,
		auth:      service,
		scheduler: innerScheduler,
		queries:   queries,
	}

	if err := scheduler.loadJobs(ctx); err != nil {
		return nil, fmt.Errorf("failed to load jobs: %w", err)
	}

	return scheduler, nil
}

func (s *Scheduler) AddReaction(reaction *sqlc.Reaction) error {
	var schedule Schedule
	if err := json.Unmarshal([]byte(reaction.Triggerdata), &schedule); err != nil {
		return fmt.Errorf("failed to unmarshal schedule data: %w", err)
	}

	_, err := s.scheduler.NewJob(
		gocron.CronJob(schedule.Expression, false),
		gocron.NewTask(
			func(ctx context.Context) {
				_, err := action.Run(ctx, s.config, s.auth, s.queries, reaction.Action, reaction.Actiondata, "{}")
				if err != nil {
					slog.ErrorContext(ctx, "Failed to run schedule reaction", "error", err, "reaction_id", reaction.ID)
				}
			},
		),
		gocron.WithTags(reaction.ID),
	)
	if err != nil {
		return fmt.Errorf("failed to create new job for reaction %s: %w", reaction.ID, err)
	}

	return nil
}

func (s *Scheduler) RemoveReaction(id string) {
	s.scheduler.RemoveByTags(id)
}

func (s *Scheduler) loadJobs(ctx context.Context) error {
	reactions, err := s.queries.ListReactions(ctx, sqlc.ListReactionsParams{
		Offset: 0,
		Limit:  1000,
	})
	if err != nil {
		return fmt.Errorf("failed to find schedule reaction: %w", err)
	}

	if len(reactions) == 0 {
		return nil
	}

	var errs []error

	for _, reaction := range reactions {
		if reaction.Trigger != "schedule" {
			continue
		}

		if err := s.AddReaction(&sqlc.Reaction{
			Action:      reaction.Action,
			Actiondata:  reaction.Actiondata,
			Created:     reaction.Created,
			ID:          reaction.ID,
			Name:        reaction.Name,
			Trigger:     reaction.Trigger,
			Triggerdata: reaction.Triggerdata,
			Updated:     reaction.Updated,
		}); err != nil {
			errs = append(errs, fmt.Errorf("failed to add reaction %s: %w", reaction.ID, err))
		}
	}

	return errors.Join(errs...)
}
