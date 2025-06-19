CREATE TABLE IF NOT EXISTS _migrations
(
    file    VARCHAR(255) PRIMARY KEY NOT NULL,
    applied INTEGER                  NOT NULL
);
CREATE TABLE IF NOT EXISTS _collections
(
    id         TEXT PRIMARY KEY                                 NOT NULL,
    system     BOOLEAN DEFAULT FALSE                            NOT NULL,
    type       TEXT    DEFAULT "base"                           NOT NULL,
    name       TEXT UNIQUE                                      NOT NULL,
    schema     JSON    DEFAULT "[]"                             NOT NULL,
    indexes    JSON    DEFAULT "[]"                             NOT NULL,
    listRule   TEXT    DEFAULT NULL,
    viewRule   TEXT    DEFAULT NULL,
    createRule TEXT    DEFAULT NULL,
    updateRule TEXT    DEFAULT NULL,
    deleteRule TEXT    DEFAULT NULL,
    options    JSON    DEFAULT "{}"                             NOT NULL,
    created    TEXT    DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL,
    updated    TEXT    DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL
);
CREATE TABLE IF NOT EXISTS _externalAuths
(
    id           TEXT PRIMARY KEY                              NOT NULL,
    collectionId TEXT                                          NOT NULL,
    recordId     TEXT                                          NOT NULL,
    provider     TEXT                                          NOT NULL,
    providerId   TEXT                                          NOT NULL,
    created      TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL,
    updated      TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL,
    ---
    FOREIGN KEY (collectionId) REFERENCES _collections (id) ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS _externalAuths_record_provider_idx on _externalAuths (collectionId, recordId, provider);
CREATE UNIQUE INDEX IF NOT EXISTS _externalAuths_collection_provider_idx on _externalAuths (collectionId, provider, providerId);
CREATE TABLE IF NOT EXISTS _admins
(
    id              TEXT PRIMARY KEY                                 NOT NULL,
    avatar          INTEGER DEFAULT 0                                NOT NULL,
    email           TEXT UNIQUE                                      NOT NULL,
    tokenKey        TEXT UNIQUE                                      NOT NULL,
    passwordHash    TEXT                                             NOT NULL,
    lastResetSentAt TEXT    DEFAULT ""                               NOT NULL,
    created         TEXT    DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL,
    updated         TEXT    DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL
);