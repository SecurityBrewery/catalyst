package database

import (
	"context"

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
