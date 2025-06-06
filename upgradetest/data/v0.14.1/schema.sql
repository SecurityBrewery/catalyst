CREATE TABLE `_migrations` (file VARCHAR(255) PRIMARY KEY NOT NULL, applied INTEGER NOT NULL);
CREATE TABLE `_admins` (
				`id`              TEXT PRIMARY KEY NOT NULL,
				`avatar`          INTEGER DEFAULT 0 NOT NULL,
				`email`           TEXT UNIQUE NOT NULL,
				`tokenKey`        TEXT UNIQUE NOT NULL,
				`passwordHash`    TEXT NOT NULL,
				`lastResetSentAt` TEXT DEFAULT "" NOT NULL,
				`created`         TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL,
				`updated`         TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL
			);
CREATE TABLE `_collections` (
				`id`         TEXT PRIMARY KEY NOT NULL,
				`system`     BOOLEAN DEFAULT FALSE NOT NULL,
				`type`       TEXT DEFAULT "base" NOT NULL,
				`name`       TEXT UNIQUE NOT NULL,
				`schema`     JSON DEFAULT "[]" NOT NULL,
				`indexes`    JSON DEFAULT "[]" NOT NULL,
				`listRule`   TEXT DEFAULT NULL,
				`viewRule`   TEXT DEFAULT NULL,
				`createRule` TEXT DEFAULT NULL,
				`updateRule` TEXT DEFAULT NULL,
				`deleteRule` TEXT DEFAULT NULL,
				`options`    JSON DEFAULT "{}" NOT NULL,
				`created`    TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL,
				`updated`    TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL
			);
CREATE TABLE `_params` (
				`id`      TEXT PRIMARY KEY NOT NULL,
				`key`     TEXT UNIQUE NOT NULL,
				`value`   JSON DEFAULT NULL,
				`created` TEXT DEFAULT "" NOT NULL,
				`updated` TEXT DEFAULT "" NOT NULL
			);
CREATE TABLE `_externalAuths` (
				`id`           TEXT PRIMARY KEY NOT NULL,
				`collectionId` TEXT NOT NULL,
				`recordId`     TEXT NOT NULL,
				`provider`     TEXT NOT NULL,
				`providerId`   TEXT NOT NULL,
				`created`      TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL,
				`updated`      TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL,
				---
				FOREIGN KEY (`collectionId`) REFERENCES `_collections` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
			);
CREATE UNIQUE INDEX _externalAuths_record_provider_idx on `_externalAuths` (`collectionId`, `recordId`, `provider`);
CREATE UNIQUE INDEX _externalAuths_collection_provider_idx on `_externalAuths` (`collectionId`, `provider`, `providerId`);
CREATE TABLE `users` (`avatar` TEXT DEFAULT '' NOT NULL, `created` TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL, `email` TEXT DEFAULT '' NOT NULL, `emailVisibility` BOOLEAN DEFAULT FALSE NOT NULL, `id` TEXT PRIMARY KEY DEFAULT ('r'||lower(hex(randomblob(7)))) NOT NULL, `lastLoginAlertSentAt` TEXT DEFAULT '' NOT NULL, `lastResetSentAt` TEXT DEFAULT '' NOT NULL, `lastVerificationSentAt` TEXT DEFAULT '' NOT NULL, `name` TEXT DEFAULT '' NOT NULL, `passwordHash` TEXT NOT NULL, `tokenKey` TEXT NOT NULL, `updated` TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL, `username` TEXT NOT NULL, `verified` BOOLEAN DEFAULT FALSE NOT NULL);
CREATE UNIQUE INDEX __pb_users_auth__username_idx ON `users` (`username`);
CREATE UNIQUE INDEX __pb_users_auth__email_idx ON `users` (`email`) WHERE `email` != '';
CREATE UNIQUE INDEX __pb_users_auth__tokenKey_idx ON `users` (`tokenKey`);
CREATE TABLE `webhooks` (`collection` TEXT DEFAULT '' NOT NULL, `created` TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL, `destination` TEXT DEFAULT '' NOT NULL, `id` TEXT PRIMARY KEY DEFAULT ('r'||lower(hex(randomblob(7)))) NOT NULL, `name` TEXT DEFAULT '' NOT NULL, `updated` TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL);
CREATE TABLE `types` (`created` TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL, `icon` TEXT DEFAULT '' NOT NULL, `id` TEXT PRIMARY KEY DEFAULT ('r'||lower(hex(randomblob(7)))) NOT NULL, `plural` TEXT DEFAULT '' NOT NULL, `schema` JSON DEFAULT NULL, `singular` TEXT DEFAULT '' NOT NULL, `updated` TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL);
CREATE TABLE `tickets` (`created` TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL, `description` TEXT DEFAULT '' NOT NULL, `id` TEXT PRIMARY KEY DEFAULT ('r'||lower(hex(randomblob(7)))) NOT NULL, `name` TEXT DEFAULT '' NOT NULL, `open` BOOLEAN DEFAULT FALSE NOT NULL, `owner` TEXT DEFAULT '' NOT NULL, `resolution` TEXT DEFAULT '' NOT NULL, `schema` JSON DEFAULT NULL, `state` JSON DEFAULT NULL, `type` TEXT DEFAULT '' NOT NULL, `updated` TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL);
CREATE TABLE `tasks` (`created` TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL, `id` TEXT PRIMARY KEY DEFAULT ('r'||lower(hex(randomblob(7)))) NOT NULL, `name` TEXT DEFAULT '' NOT NULL, `open` BOOLEAN DEFAULT FALSE NOT NULL, `owner` TEXT DEFAULT '' NOT NULL, `ticket` TEXT DEFAULT '' NOT NULL, `updated` TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL);
CREATE TABLE `comments` (`author` TEXT DEFAULT '' NOT NULL, `created` TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL, `id` TEXT PRIMARY KEY DEFAULT ('r'||lower(hex(randomblob(7)))) NOT NULL, `message` TEXT DEFAULT '' NOT NULL, `ticket` TEXT DEFAULT '' NOT NULL, `updated` TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL);
CREATE TABLE `timeline` (`created` TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL, `id` TEXT PRIMARY KEY DEFAULT ('r'||lower(hex(randomblob(7)))) NOT NULL, `message` TEXT DEFAULT '' NOT NULL, `ticket` TEXT DEFAULT '' NOT NULL, `time` TEXT DEFAULT '' NOT NULL, `updated` TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL);
CREATE TABLE `links` (`created` TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL, `id` TEXT PRIMARY KEY DEFAULT ('r'||lower(hex(randomblob(7)))) NOT NULL, `name` TEXT DEFAULT '' NOT NULL, `ticket` TEXT DEFAULT '' NOT NULL, `updated` TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL, `url` TEXT DEFAULT '' NOT NULL);
CREATE TABLE `files` (`blob` TEXT DEFAULT '' NOT NULL, `created` TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL, `id` TEXT PRIMARY KEY DEFAULT ('r'||lower(hex(randomblob(7)))) NOT NULL, `name` TEXT DEFAULT '' NOT NULL, `size` NUMERIC DEFAULT 0 NOT NULL, `ticket` TEXT DEFAULT '' NOT NULL, `updated` TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL);
CREATE TABLE `features` (`created` TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL, `id` TEXT PRIMARY KEY DEFAULT ('r'||lower(hex(randomblob(7)))) NOT NULL, `name` TEXT DEFAULT '' NOT NULL, `updated` TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL);
CREATE UNIQUE INDEX `unique_name` ON `features` (`name`);
CREATE VIEW `sidebar` AS SELECT * FROM (SELECT types.id as id, types.singular as singular, types.plural as plural, types.icon as icon, (SELECT COUNT(tickets.id) FROM tickets WHERE tickets.type = types.id AND tickets.open = true) as count
FROM types
ORDER BY types.plural)
/* sidebar(id,singular,plural,icon,count) */;
CREATE TABLE `reactions` (`action` TEXT DEFAULT '' NOT NULL, `actiondata` JSON DEFAULT NULL, `created` TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL, `id` TEXT PRIMARY KEY DEFAULT ('r'||lower(hex(randomblob(7)))) NOT NULL, `name` TEXT DEFAULT '' NOT NULL, `trigger` TEXT DEFAULT '' NOT NULL, `triggerdata` JSON DEFAULT NULL, `updated` TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL);
CREATE VIEW `ticket_search` AS SELECT * FROM (SELECT 
    tickets.id, 
    tickets.name,
    tickets.created,
    tickets.description,
    tickets.open,
    tickets.type,
    tickets.state,
    users.name as owner_name,
    group_concat(comments.message) as comment_messages,
    group_concat(files.name) as file_names,
    group_concat(links.name) as link_names,
    group_concat(links.url) as link_urls,
    group_concat(tasks.name) as task_names,
    group_concat(timeline.message) as timeline_messages
FROM tickets
LEFT JOIN comments ON comments.ticket = tickets.id
LEFT JOIN files ON files.ticket = tickets.id
LEFT JOIN links ON links.ticket = tickets.id
LEFT JOIN tasks ON tasks.ticket = tickets.id
LEFT JOIN timeline ON timeline.ticket = tickets.id
LEFT JOIN users ON users.id = tickets.owner
GROUP BY tickets.id)
/* ticket_search(id,name,created,description,open,type,state,owner_name,comment_messages,file_names,link_names,link_urls,task_names,timeline_messages) */;
CREATE VIEW `dashboard_counts` AS SELECT * FROM (SELECT cast(`id` as text) `id`,`count` FROM (SELECT id, count FROM (
  SELECT 'users' as id, COUNT(users.id) as count FROM users
  UNION
  SELECT 'tickets' as id, COUNT(tickets.id) as count FROM tickets
  UNION
  SELECT 'tasks' as id, COUNT(tasks.id) as count FROM tasks
  UNION 
  SELECT 'reactions' as id, COUNT(reactions.id) as count FROM reactions
) as counts))
/* dashboard_counts(id,count) */;
