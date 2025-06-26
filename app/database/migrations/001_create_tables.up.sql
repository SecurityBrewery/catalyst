DROP TABLE _migrations;
DROP TABLE _collections;
DROP TABLE _externalauths;
DROP VIEW sidebar;
DROP VIEW ticket_search;
DROP VIEW dashboard_counts;

--- _params

CREATE TABLE new_params
(
    key   TEXT PRIMARY KEY NOT NULL,
    value JSON
);

INSERT INTO new_params
    (key, value)
SELECT key, value
FROM _params;

DROP TABLE _params;
ALTER TABLE new_params
    RENAME TO _params;

--- users

CREATE TABLE new_users
(
    id                     TEXT PRIMARY KEY DEFAULT ('u' || lower(hex(randomblob(7)))) NOT NULL,
    username               TEXT                                                        NOT NULL,
    passwordHash           TEXT                                                        NOT NULL,
    tokenKey               TEXT                                                        NOT NULL,
    active                 BOOLEAN                                                     NOT NULL,
    name                   TEXT,
    email                  TEXT,
    avatar                 TEXT,
    lastresetsentat        DATETIME,
    lastverificationsentat DATETIME,
    admin                  BOOLEAN                                                     NOT NULL,
    created                DATETIME         DEFAULT CURRENT_TIMESTAMP                  NOT NULL,
    updated                DATETIME         DEFAULT CURRENT_TIMESTAMP                  NOT NULL
);

INSERT INTO new_users
(avatar, email, id, lastresetsentat, lastverificationsentat, name, passwordHash, tokenKey, username, active, admin,
 created,
 updated)
SELECT avatar,
       email,
       id,
       lastResetSentAt,
       lastVerificationSentAt,
       name,
       passwordHash,
       tokenKey,
       username,
       verified,
       false,
       created,
       updated
FROM users;

INSERT INTO new_users
(avatar, email, id, lastresetsentat, lastverificationsentat, name, passwordHash, tokenKey, username, active, admin,
 created,
 updated)
SELECT avatar,
       email,
       id,
       lastResetSentAt,
       '',
       email,
       passwordHash,
       tokenKey,
       id,
       true,
       true,
       created,
       updated
FROM _admins;

DROP TABLE users;
DROP TABLE _admins;
ALTER TABLE new_users
    RENAME TO users;

--- webhooks

CREATE TABLE new_webhooks
(
    id          TEXT PRIMARY KEY DEFAULT ('w' || lower(hex(randomblob(7)))) NOT NULL,
    collection  TEXT                                                        NOT NULL,
    destination TEXT                                                        NOT NULL,
    name        TEXT                                                        NOT NULL,
    created     DATETIME         DEFAULT CURRENT_TIMESTAMP                  NOT NULL,
    updated     DATETIME         DEFAULT CURRENT_TIMESTAMP                  NOT NULL
);

INSERT INTO new_webhooks
    (collection, destination, id, name, created, updated)
SELECT collection, destination, id, name, datetime(created), datetime(updated)
FROM webhooks;

DROP TABLE webhooks;
ALTER TABLE new_webhooks
    RENAME TO webhooks;

--- types

CREATE TABLE new_types
(
    id       TEXT PRIMARY KEY DEFAULT ('y' || lower(hex(randomblob(7)))) NOT NULL,
    icon     TEXT,
    singular TEXT                                                        NOT NULL,
    plural   TEXT                                                        NOT NULL,
    schema   JSON,
    created  DATETIME         DEFAULT CURRENT_TIMESTAMP                  NOT NULL,
    updated  DATETIME         DEFAULT CURRENT_TIMESTAMP                  NOT NULL
);

INSERT INTO new_types
    (id, icon, singular, plural, schema, created, updated)
SELECT id, icon, singular, plural, schema, created, updated
FROM types;

DROP TABLE types;
ALTER TABLE new_types
    RENAME TO types;

--- ticket

CREATE TABLE new_tickets
(
    id          TEXT PRIMARY KEY DEFAULT ('t' || lower(hex(randomblob(7)))) NOT NULL,
    type        TEXT                                                        NOT NULL,
    owner       TEXT,
    name        TEXT                                                        NOT NULL,
    description TEXT                                                        NOT NULL,
    open        BOOLEAN                                                     NOT NULL,
    resolution  TEXT,
    schema      JSON,
    state       JSON,
    created     DATETIME         DEFAULT CURRENT_TIMESTAMP                  NOT NULL,
    updated     DATETIME         DEFAULT CURRENT_TIMESTAMP                  NOT NULL,

    FOREIGN KEY (type) REFERENCES types (id) ON DELETE SET NULL,
    FOREIGN KEY (owner) REFERENCES users (id) ON DELETE SET NULL
);

INSERT INTO new_tickets
(id, name, description, open, owner, resolution, schema, state, type, created, updated)
SELECT id,
       name,
       description,
       open,
       owner,
       resolution,
       schema,
       state,
       type,
       created,
       updated
FROM tickets;

DROP TABLE tickets;
ALTER TABLE new_tickets
    RENAME TO tickets;

--- tasks

CREATE TABLE new_tasks
(
    id      TEXT PRIMARY KEY DEFAULT ('t' || lower(hex(randomblob(7)))) NOT NULL,
    ticket  TEXT                                                        NOT NULL,
    owner   TEXT,
    name    TEXT                                                        NOT NULL,
    open    BOOLEAN                                                     NOT NULL,
    created DATETIME         DEFAULT CURRENT_TIMESTAMP                  NOT NULL,
    updated DATETIME         DEFAULT CURRENT_TIMESTAMP                  NOT NULL,

    FOREIGN KEY (ticket) REFERENCES tickets (id) ON DELETE CASCADE,
    FOREIGN KEY (owner) REFERENCES users (id) ON DELETE SET NULL
);

INSERT INTO new_tasks
    (id, ticket, owner, name, open, created, updated)
SELECT id, ticket, owner, name, open, created, updated
FROM tasks;
DROP TABLE tasks;
ALTER TABLE new_tasks
    RENAME TO tasks;

--- comments

CREATE TABLE new_comments
(
    id      TEXT PRIMARY KEY DEFAULT ('c' || lower(hex(randomblob(7)))) NOT NULL,
    ticket  TEXT                                                        NOT NULL,
    author  TEXT                                                        NOT NULL,
    message TEXT                                                        NOT NULL,
    created DATETIME         DEFAULT CURRENT_TIMESTAMP                  NOT NULL,
    updated DATETIME         DEFAULT CURRENT_TIMESTAMP                  NOT NULL,

    FOREIGN KEY (ticket) REFERENCES tickets (id) ON DELETE CASCADE,
    FOREIGN KEY (author) REFERENCES users (id) ON DELETE SET NULL
);

INSERT INTO new_comments
    (id, ticket, author, message, created, updated)
SELECT id, ticket, author, message, created, updated
FROM comments;
DROP TABLE comments;
ALTER TABLE new_comments
    RENAME TO comments;

--- timeline

CREATE TABLE new_timeline
(
    id      TEXT PRIMARY KEY DEFAULT ('h' || lower(hex(randomblob(7)))) NOT NULL,
    ticket  TEXT                                                        NOT NULL,
    message TEXT                                                        NOT NULL,
    time    DATETIME                                                    NOT NULL,
    created DATETIME         DEFAULT CURRENT_TIMESTAMP                  NOT NULL,
    updated DATETIME         DEFAULT CURRENT_TIMESTAMP                  NOT NULL,

    FOREIGN KEY (ticket) REFERENCES tickets (id) ON DELETE CASCADE
);

INSERT INTO new_timeline
    (id, ticket, message, time, created, updated)
SELECT id, ticket, message, time, created, updated
FROM timeline;

DROP TABLE timeline;
ALTER TABLE new_timeline
    RENAME TO timeline;

--- links

CREATE TABLE new_links
(
    id      TEXT PRIMARY KEY DEFAULT ('l' || lower(hex(randomblob(7)))) NOT NULL,
    ticket  TEXT                                                        NOT NULL,
    name    TEXT                                                        NOT NULL,
    url     TEXT                                                        NOT NULL,
    created DATETIME         DEFAULT CURRENT_TIMESTAMP                  NOT NULL,
    updated DATETIME         DEFAULT CURRENT_TIMESTAMP                  NOT NULL,

    FOREIGN KEY (ticket) REFERENCES tickets (id) ON DELETE CASCADE
);

INSERT INTO new_links
    (id, ticket, name, url, created, updated)
SELECT id, ticket, name, url, datetime(created), datetime(updated)
FROM links;
DROP TABLE links;
ALTER TABLE new_links
    RENAME TO links;

--- files

CREATE TABLE new_files
(
    id      TEXT PRIMARY KEY DEFAULT ('b' || lower(hex(randomblob(7)))) NOT NULL,
    ticket  TEXT                                                        NOT NULL,
    name    TEXT                                                        NOT NULL,
    blob    TEXT                                                        NOT NULL,
    size    NUMERIC                                                     NOT NULL,
    created DATETIME         DEFAULT CURRENT_TIMESTAMP                  NOT NULL,
    updated DATETIME         DEFAULT CURRENT_TIMESTAMP                  NOT NULL,

    FOREIGN KEY (ticket) REFERENCES tickets (id) ON DELETE CASCADE
);

INSERT INTO new_files
    (id, name, blob, size, ticket, created, updated)
SELECT id, name, blob, size, ticket, created, updated
FROM files;
DROP TABLE files;
ALTER TABLE new_files
    RENAME TO files;

--- features

CREATE TABLE new_features
(
    key TEXT PRIMARY KEY NOT NULL
);

INSERT INTO new_features
    (key)
SELECT name
FROM features;

DROP TABLE features;
ALTER TABLE new_features
    RENAME TO features;

--- reactions

CREATE TABLE new_reactions
(
    id          TEXT PRIMARY KEY DEFAULT ('r' || lower(hex(randomblob(7)))) NOT NULL,
    name        TEXT                                                        NOT NULL,
    action      TEXT                                                        NOT NULL,
    actiondata  JSON                                                        NOT NULL,
    trigger     TEXT                                                        NOT NULL,
    triggerdata JSON                                                        NOT NULL,
    created     DATETIME         DEFAULT CURRENT_TIMESTAMP                  NOT NULL,
    updated     DATETIME         DEFAULT CURRENT_TIMESTAMP                  NOT NULL
);

INSERT INTO new_reactions
    (id, name, action, actiondata, trigger, triggerdata, created, updated)
SELECT id,
       name,
       action,
       actionData,
       trigger,
       triggerData,
       created,
       updated
FROM reactions;
DROP TABLE reactions;
ALTER TABLE new_reactions
    RENAME TO reactions;

--- views

CREATE VIEW sidebar AS
SELECT types.id                                                                                      as id,
       types.singular                                                                                as singular,
       types.plural                                                                                  as plural,
       types.icon                                                                                    as icon,
       (SELECT COUNT(tickets.id) FROM tickets WHERE tickets.type = types.id AND tickets.open = true) as count
FROM types
ORDER BY types.plural;

CREATE VIEW ticket_search AS
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

CREATE VIEW dashboard_counts AS
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