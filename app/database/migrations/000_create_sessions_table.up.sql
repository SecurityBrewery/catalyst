CREATE TABLE IF NOT EXISTS sessions
(
    token  TEXT PRIMARY KEY,
    data   BLOB    NOT NULL,
    expiry INTEGER NOT NULL
) STRICT;

CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions (expiry);
