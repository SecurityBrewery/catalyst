DROP VIEW dashboard_counts;
DROP VIEW ticket_search;
DROP VIEW sidebar;

DROP TABLE reactions;
DROP TABLE features;
DROP TABLE files;
DROP TABLE links;
DROP TABLE timeline;
DROP TABLE comments;
DROP TABLE tasks;
DROP TABLE tickets;
DROP TABLE types;
DROP TABLE webhooks;
DROP TABLE users;
DROP TABLE _params;

CREATE TABLE IF NOT EXISTS _migrations
(
    file    VARCHAR(255) PRIMARY KEY NOT NULL,
    applied INTEGER                  NOT NULL
);
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
CREATE TABLE IF NOT EXISTS _params
(
    id      TEXT PRIMARY KEY NOT NULL,
    key     TEXT UNIQUE      NOT NULL,
    value   JSON DEFAULT NULL,
    created TEXT DEFAULT ""  NOT NULL,
    updated TEXT DEFAULT ""  NOT NULL
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
CREATE TABLE IF NOT EXISTS users
(
    avatar                 TEXT             DEFAULT ''                                 NOT NULL,
    created                TEXT             DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ'))   NOT NULL,
    email                  TEXT             DEFAULT ''                                 NOT NULL,
    emailVisibility        BOOLEAN          DEFAULT FALSE                              NOT NULL,
    id                     TEXT PRIMARY KEY DEFAULT ('r' || lower(hex(randomblob(7)))) NOT NULL,
    lastLoginAlertSentAt   TEXT             DEFAULT ''                                 NOT NULL,
    lastResetSentAt        TEXT             DEFAULT ''                                 NOT NULL,
    lastVerificationSentAt TEXT             DEFAULT ''                                 NOT NULL,
    name                   TEXT             DEFAULT ''                                 NOT NULL,
    passwordHash           TEXT                                                        NOT NULL,
    tokenKey               TEXT                                                        NOT NULL,
    updated                TEXT             DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ'))   NOT NULL,
    username               TEXT                                                        NOT NULL,
    verified               BOOLEAN          DEFAULT FALSE                              NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS __pb_users_auth__username_idx ON users (username);
CREATE UNIQUE INDEX IF NOT EXISTS __pb_users_auth__email_idx ON users (email) WHERE email != '';
CREATE UNIQUE INDEX IF NOT EXISTS __pb_users_auth__tokenKey_idx ON users (tokenKey);
CREATE TABLE IF NOT EXISTS webhooks
(
    collection  TEXT             DEFAULT ''                                 NOT NULL,
    created     TEXT             DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ'))   NOT NULL,
    destination TEXT             DEFAULT ''                                 NOT NULL,
    id          TEXT PRIMARY KEY DEFAULT ('r' || lower(hex(randomblob(7)))) NOT NULL,
    name        TEXT             DEFAULT ''                                 NOT NULL,
    updated     TEXT             DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ'))   NOT NULL
);
CREATE TABLE IF NOT EXISTS types
(
    created  TEXT             DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ'))   NOT NULL,
    icon     TEXT             DEFAULT ''                                 NOT NULL,
    id       TEXT PRIMARY KEY DEFAULT ('r' || lower(hex(randomblob(7)))) NOT NULL,
    plural   TEXT             DEFAULT ''                                 NOT NULL,
    schema   JSON             DEFAULT NULL,
    singular TEXT             DEFAULT ''                                 NOT NULL,
    updated  TEXT             DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ'))   NOT NULL
);
CREATE TABLE IF NOT EXISTS tickets
(
    created     TEXT             DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ'))   NOT NULL,
    description TEXT             DEFAULT ''                                 NOT NULL,
    id          TEXT PRIMARY KEY DEFAULT ('r' || lower(hex(randomblob(7)))) NOT NULL,
    name        TEXT             DEFAULT ''                                 NOT NULL,
    open        BOOLEAN          DEFAULT FALSE                              NOT NULL,
    owner       TEXT             DEFAULT ''                                 NOT NULL,
    resolution  TEXT             DEFAULT ''                                 NOT NULL,
    schema      JSON             DEFAULT NULL,
    state       JSON             DEFAULT NULL,
    type        TEXT             DEFAULT ''                                 NOT NULL,
    updated     TEXT             DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ'))   NOT NULL
);
CREATE TABLE IF NOT EXISTS tasks
(
    created TEXT             DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ'))   NOT NULL,
    id      TEXT PRIMARY KEY DEFAULT ('r' || lower(hex(randomblob(7)))) NOT NULL,
    name    TEXT             DEFAULT ''                                 NOT NULL,
    open    BOOLEAN          DEFAULT FALSE                              NOT NULL,
    owner   TEXT             DEFAULT ''                                 NOT NULL,
    ticket  TEXT             DEFAULT ''                                 NOT NULL,
    updated TEXT             DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ'))   NOT NULL
);
CREATE TABLE IF NOT EXISTS comments
(
    author  TEXT             DEFAULT ''                                 NOT NULL,
    created TEXT             DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ'))   NOT NULL,
    id      TEXT PRIMARY KEY DEFAULT ('r' || lower(hex(randomblob(7)))) NOT NULL,
    message TEXT             DEFAULT ''                                 NOT NULL,
    ticket  TEXT             DEFAULT ''                                 NOT NULL,
    updated TEXT             DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ'))   NOT NULL
);
CREATE TABLE IF NOT EXISTS timeline
(
    created TEXT             DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ'))   NOT NULL,
    id      TEXT PRIMARY KEY DEFAULT ('r' || lower(hex(randomblob(7)))) NOT NULL,
    message TEXT             DEFAULT ''                                 NOT NULL,
    ticket  TEXT             DEFAULT ''                                 NOT NULL,
    time    TEXT             DEFAULT ''                                 NOT NULL,
    updated TEXT             DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ'))   NOT NULL
);
CREATE TABLE IF NOT EXISTS links
(
    created TEXT             DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ'))   NOT NULL,
    id      TEXT PRIMARY KEY DEFAULT ('r' || lower(hex(randomblob(7)))) NOT NULL,
    name    TEXT             DEFAULT ''                                 NOT NULL,
    ticket  TEXT             DEFAULT ''                                 NOT NULL,
    updated TEXT             DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ'))   NOT NULL,
    url     TEXT             DEFAULT ''                                 NOT NULL
);
CREATE TABLE IF NOT EXISTS files
(
    blob    TEXT             DEFAULT ''                                 NOT NULL,
    created TEXT             DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ'))   NOT NULL,
    id      TEXT PRIMARY KEY DEFAULT ('r' || lower(hex(randomblob(7)))) NOT NULL,
    name    TEXT             DEFAULT ''                                 NOT NULL,
    size    NUMERIC          DEFAULT 0                                  NOT NULL,
    ticket  TEXT             DEFAULT ''                                 NOT NULL,
    updated TEXT             DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ'))   NOT NULL
);
CREATE TABLE IF NOT EXISTS features
(
    created TEXT             DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ'))   NOT NULL,
    id      TEXT PRIMARY KEY DEFAULT ('r' || lower(hex(randomblob(7)))) NOT NULL,
    name    TEXT             DEFAULT ''                                 NOT NULL,
    updated TEXT             DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ'))   NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS unique_name ON features (name);

CREATE VIEW IF NOT EXISTS sidebar AS
SELECT types.id                                                                                      as id,
       types.singular                                                                                as singular,
       types.plural                                                                                  as plural,
       types.icon                                                                                    as icon,
       (SELECT COUNT(tickets.id) FROM tickets WHERE tickets.type = types.id AND tickets.open = true) as count
FROM types
ORDER BY types.plural;

CREATE TABLE IF NOT EXISTS reactions
(
    action      TEXT             DEFAULT ''                                 NOT NULL,
    actiondata  JSON             DEFAULT NULL,
    created     TEXT             DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ'))   NOT NULL,
    id          TEXT PRIMARY KEY DEFAULT ('r' || lower(hex(randomblob(7)))) NOT NULL,
    name        TEXT             DEFAULT ''                                 NOT NULL,
    trigger     TEXT             DEFAULT ''                                 NOT NULL,
    triggerdata JSON             DEFAULT NULL,
    updated     TEXT             DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ'))   NOT NULL
);

CREATE VIEW IF NOT EXISTS ticket_search AS
SELECT tickets.id,
       tickets.name,
       tickets.created,
       tickets.description,
       tickets.open,
       tickets.type,
       tickets.state,
       users.name as owner_name,
       group_concat(comments.message
       )          as comment_messages,
       group_concat(files.name
       )          as file_names,
       group_concat(links.name
       )          as link_names,
       group_concat(links.url
       )          as link_urls,
       group_concat(tasks.name
       )          as task_names,
       group_concat(timeline.message
       )          as timeline_messages
FROM tickets
         LEFT JOIN comments ON comments.ticket = tickets.id
         LEFT JOIN files ON files.ticket = tickets.id
         LEFT JOIN links ON links.ticket = tickets.id
         LEFT JOIN tasks ON tasks.ticket = tickets.id
         LEFT JOIN timeline ON timeline.ticket = tickets.id
         LEFT JOIN users ON users.id = tickets.owner
GROUP BY tickets.id;

CREATE VIEW IF NOT EXISTS dashboard_counts AS
SELECT id, count
FROM (SELECT 'users' as id,
             COUNT(users.id
             )       as count
      FROM users
      UNION
      SELECT 'tickets' as id,
             COUNT(tickets.id
             )         as count
      FROM tickets
      UNION
      SELECT 'tasks' as id,
             COUNT(tasks.id
             )       as count
      FROM tasks
      UNION
      SELECT 'reactions' as id,
             COUNT(reactions.id
             )           as count
      FROM reactions) as counts;