package sessionmanager

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/alexedwards/scs/v2"

	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

const SessionKey = "user_id"

type Config struct {
	Domain       string `json:"domain" yaml:"domain"`
	CookieSecure bool   `json:"cookieSecure,omitempty" yaml:"cookieSecure,omitempty"`
}
type SessionManager struct {
	queries  *sqlc.Queries
	internal *scs.SessionManager
}

func New(config *Config, queries *sqlc.Queries) *SessionManager {
	sessionManager := scs.New()
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Domain = config.Domain
	sessionManager.Cookie.Secure = config.CookieSecure
	sessionManager.Codec = &JSONCodec{}
	sessionManager.Store = &SQliteStore{
		queries: queries,
	}

	return &SessionManager{
		queries:  queries,
		internal: sessionManager,
	}
}

func (m *SessionManager) Get(ctx context.Context) (sqlc.User, error, error) {
	v := m.internal.Get(ctx, SessionKey)

	id, ok := v.(string)
	if !ok {
		return sqlc.User{}, errors.New("no session found"), nil
	}

	user, err := m.queries.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, errors.New("user not found"), nil
		}

		return sqlc.User{}, nil, err
	}

	return user, nil, nil
}

func (m *SessionManager) Exists(ctx context.Context) bool {
	return m.internal.Exists(ctx, SessionKey)
}

func (m *SessionManager) Put(ctx context.Context, id string) {
	m.internal.Put(ctx, SessionKey, id)
}

func (m *SessionManager) Remove(ctx context.Context) {
	m.internal.Remove(ctx, SessionKey)
}

func (m *SessionManager) LoadAndSave(handler http.Handler) http.Handler {
	return m.internal.LoadAndSave(handler)
}
