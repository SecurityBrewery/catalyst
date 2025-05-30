package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/alexedwards/scs/v2"

	"github.com/SecurityBrewery/catalyst/app2/database/sqlc"
)

var _ scs.CtxStore = &SQliteStore{}

type SQliteStore struct {
	queries *sqlc.Queries
}

func (f *SQliteStore) Find(string) ([]byte, bool, error)      { panic("implement me") }
func (f *SQliteStore) Commit(string, []byte, time.Time) error { panic("implement me") }
func (f *SQliteStore) Delete(string) error                    { panic("implement me") }

func (f *SQliteStore) FindCtx(ctx context.Context, token string) (b []byte, found bool, err error) {
	session, err := f.queries.FindSession(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}

		return nil, false, err
	}

	return session.Data, true, nil
}

func (f *SQliteStore) CommitCtx(ctx context.Context, token string, b []byte, expiry time.Time) (err error) {
	return f.queries.CommitSession(ctx, sqlc.CommitSessionParams{
		Token:  token,
		Data:   b,
		Expiry: expiry.Unix(),
	})
}

func (f *SQliteStore) DeleteCtx(ctx context.Context, token string) (err error) {
	if err := f.queries.DeleteSession(ctx, token); err != nil {
		return err
	}

	return nil
}
