package database

import (
	"context"
	"fmt"

	"github.com/SecurityBrewery/catalyst/caql"
	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/model"
)

func (db *Database) Statistics(ctx context.Context) (*model.Statistics, error) {
	query := `RETURN { 
	tickets_per_type: MERGE(FOR d in tickets
		COLLECT type = d.type WITH COUNT INTO typecount
		RETURN ZIP([type], [typecount])), 
		
	unassigned: FIRST(FOR d in tickets
		FILTER d.status == "open" AND !d.owner
		COLLECT WITH COUNT INTO length
		RETURN length),
	
	open_tickets_per_user: MERGE(FOR d in tickets
		FILTER d.status == "open"
		COLLECT user = d.owner WITH COUNT INTO usercount
		RETURN ZIP([user], [usercount])), 
	
	tickets_per_week: MERGE(FOR d in tickets
		COLLECT week = CONCAT(DATE_YEAR(d.created), "-", DATE_ISOWEEK(d.created) < 10 ? "0" : "", DATE_ISOWEEK(d.created)) WITH COUNT INTO weekcount
		RETURN ZIP([week], [weekcount])),
	}`

	cur, _, err := db.Query(ctx, query, nil, busdb.ReadOperation)
	if err != nil {
		return nil, err
	}
	defer cur.Close()

	statistics := model.Statistics{}
	if _, err := cur.ReadDocument(ctx, &statistics); err != nil {
		return nil, err
	}

	return &statistics, nil
}

func (db *Database) WidgetData(ctx context.Context, aggregation string, filter *string) (map[string]any, error) {
	parser := &caql.Parser{Searcher: db.Index, Prefix: "d."}

	queryTree, err := parser.Parse(aggregation)
	if err != nil {
		return nil, fmt.Errorf("invalid aggregation query (%s): syntax error", aggregation)
	}
	aggregationString, err := queryTree.String()
	if err != nil {
		return nil, fmt.Errorf("invalid widget aggregation query (%s): %w", aggregation, err)
	}
	aggregation = aggregationString

	filterQ := ""
	if filter != nil && *filter != "" {
		queryTree, err := parser.Parse(*filter)
		if err != nil {
			return nil, fmt.Errorf("invalid filter query (%s): syntax error", *filter)
		}
		filterString, err := queryTree.String()
		if err != nil {
			return nil, fmt.Errorf("invalid widget filter query (%s): %w", *filter, err)
		}

		filterQ = "FILTER " + filterString
	}

	query := `RETURN MERGE(FOR d in tickets
		` + filterQ + `
		COLLECT field = ` + aggregation + ` WITH COUNT INTO count
		RETURN ZIP([field], [count]))`

	cur, _, err := db.Query(ctx, query, nil, busdb.ReadOperation)
	if err != nil {
		return nil, err
	}
	defer cur.Close()

	statistics := map[string]any{}
	if _, err := cur.ReadDocument(ctx, &statistics); err != nil {
		return nil, err
	}

	return statistics, nil
}
