CREATE TABLE IF NOT EXISTS comments
(
    author  TEXT default ''                                 not null,
    created TEXT default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null,
    id      TEXT default ('r' || lower(hex(randomblob(7)))) not null
        primary key,
    message TEXT default ''                                 not null,
    ticket  TEXT default ''                                 not null,
    updated TEXT default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null
);

CREATE TABLE IF NOT EXISTS features
(
    created TEXT default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null,
    id      TEXT default ('r' || lower(hex(randomblob(7)))) not null
        primary key,
    name    TEXT default ''                                 not null,
    updated TEXT default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null
);

create unique index if not exists unique_name
    on features (name);

CREATE TABLE IF NOT EXISTS files
(
    blob    TEXT    default ''                                 not null,
    created TEXT    default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null,
    id      TEXT    default ('r' || lower(hex(randomblob(7)))) not null
        primary key,
    name    TEXT    default ''                                 not null,
    size    INTEGER default 0                                  not null,
    ticket  TEXT    default ''                                 not null,
    updated TEXT    default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null
);

CREATE TABLE IF NOT EXISTS links
(
    created TEXT default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null,
    id      TEXT default ('r' || lower(hex(randomblob(7)))) not null
        primary key,
    name    TEXT default ''                                 not null,
    ticket  TEXT default ''                                 not null,
    updated TEXT default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null,
    url     TEXT default ''                                 not null
);

CREATE TABLE IF NOT EXISTS reactions
(
    action      TEXT default ''                                 not null,
    actiondata  JSON default NULL,
    created     TEXT default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null,
    id          TEXT default ('r' || lower(hex(randomblob(7)))) not null
        primary key,
    name        TEXT default ''                                 not null,
    trigger     TEXT default ''                                 not null,
    triggerdata JSON default NULL,
    updated     TEXT default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null
);

CREATE TABLE IF NOT EXISTS tasks
(
    created TEXT    default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null,
    id      TEXT    default ('r' || lower(hex(randomblob(7)))) not null
        primary key,
    name    TEXT    default ''                                 not null,
    open    BOOLEAN default FALSE                              not null,
    owner   TEXT    default ''                                 not null,
    ticket  TEXT    default ''                                 not null,
    updated TEXT    default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null
);

CREATE TABLE IF NOT EXISTS tickets
(
    created     TEXT    default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null,
    description TEXT    default ''                                 not null,
    id          TEXT    default ('r' || lower(hex(randomblob(7)))) not null
        primary key,
    name        TEXT    default ''                                 not null,
    open        BOOLEAN default FALSE                              not null,
    owner       TEXT    default ''                                 not null,
    resolution  TEXT    default ''                                 not null,
    schema      JSON    default NULL,
    state       JSON    default NULL,
    type        TEXT    default ''                                 not null,
    updated     TEXT    default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null
);

CREATE TABLE IF NOT EXISTS timeline
(
    created TEXT default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null,
    id      TEXT default ('r' || lower(hex(randomblob(7)))) not null
        primary key,
    message TEXT default ''                                 not null,
    ticket  TEXT default ''                                 not null,
    time    TEXT default ''                                 not null,
    updated TEXT default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null
);

CREATE TABLE IF NOT EXISTS types
(
    created  TEXT default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null,
    icon     TEXT default ''                                 not null,
    id       TEXT default ('r' || lower(hex(randomblob(7)))) not null
        primary key,
    plural   TEXT default ''                                 not null,
    schema   JSON default NULL,
    singular TEXT default ''                                 not null,
    updated  TEXT default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null
);

CREATE TABLE IF NOT EXISTS users
(
    avatar                 TEXT    default ''                                 not null,
    created                TEXT    default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null,
    email                  TEXT    default ''                                 not null,
    emailVisibility        BOOLEAN default FALSE                              not null,
    id                     TEXT    default ('r' || lower(hex(randomblob(7)))) not null
        primary key,
    lastLoginAlertSentAt   TEXT    default ''                                 not null,
    lastResetSentAt        TEXT    default ''                                 not null,
    lastVerificationSentAt TEXT    default ''                                 not null,
    name                   TEXT    default ''                                 not null,
    passwordHash           TEXT                                               not null,
    tokenKey               TEXT                                               not null,
    updated                TEXT    default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null,
    username               TEXT                                               not null,
    verified               BOOLEAN default FALSE                              not null
);

CREATE TABLE IF NOT EXISTS webhooks
(
    collection  TEXT default ''                                 not null,
    created     TEXT default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null,
    destination TEXT default ''                                 not null,
    id          TEXT default ('r' || lower(hex(randomblob(7)))) not null
        primary key,
    name        TEXT default ''                                 not null,
    updated     TEXT default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null
);

CREATE VIEW IF NOT EXISTS dashboard_counts AS SELECT id, count FROM (
    SELECT 'users' as id, COUNT(users.id
) as count FROM users
    UNION
    SELECT 'tickets' as id, COUNT(tickets.id
) as count FROM tickets
    UNION
    SELECT 'tasks' as id, COUNT(tasks.id
) as count FROM tasks
    UNION
    SELECT 'reactions' as id, COUNT(reactions.id
) as count FROM reactions
) as counts;

CREATE VIEW IF NOT EXISTS sidebar AS
SELECT types.id                                                                                      as id,
       types.singular                                                                                as singular,
       types.plural                                                                                  as plural,
       types.icon                                                                                    as icon,
       (SELECT COUNT(tickets.id) FROM tickets WHERE tickets.type = types.id AND tickets.open = true) as count
FROM types
ORDER BY types.plural;

CREATE VIEW IF NOT EXISTS ticket_search AS SELECT
    tickets.id,
    tickets.name,
    tickets.created,
    tickets.description,
    tickets.open,
    tickets.type,
    tickets.state,
    users.name as owner_name,
    group_concat(comments.message
) as comment_messages,
    group_concat(files.name
) as file_names,
    group_concat(links.name
) as link_names,
    group_concat(links.url
) as link_urls,
    group_concat(tasks.name
) as task_names,
    group_concat(timeline.message
) as timeline_messages
    FROM tickets
    LEFT JOIN comments ON comments.ticket = tickets.id
    LEFT JOIN files ON files.ticket = tickets.id
    LEFT JOIN links ON links.ticket = tickets.id
    LEFT JOIN tasks ON tasks.ticket = tickets.id
    LEFT JOIN timeline ON timeline.ticket = tickets.id
    LEFT JOIN users ON users.id = tickets.owner
    GROUP BY tickets.id;
