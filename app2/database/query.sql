-- name: CreateTicket :execlastid
INSERT INTO tickets (name, description, open, owner, resolution, schema, state, type)
VALUES (@name, @description, @open, @owner, @resolution, @schema, @state, @type);

-- name: Ticket :one
SELECT *
FROM tickets
WHERE id = @id;

-- name: UpdateTicket :exec
UPDATE tickets
SET name = @name,
    description = @description,
    open = @open,
    owner = @owner,
    resolution = @resolution,
    schema = @schema,
    state = @state,
    type = @type
WHERE id = @id;

-- name: DeleteTicket :exec
DELETE
FROM tickets
WHERE id = @id;

-- name: ListTickets :many
SELECT *
FROM tickets
ORDER BY created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: FindSession :one
SELECT *
FROM sessions
WHERE token = @token;

-- name: CommitSession :exec
INSERT OR
REPLACE
INTO sessions (token, data, expiry)
VALUES (@token, @data, @expiry);

-- name: DeleteSession :exec
DELETE
FROM sessions
WHERE token = @token;

------------------------------------------------------------------

-- name: CreateComment :execlastid
INSERT INTO comments (author, message, ticket)
VALUES (@author, @message, @ticket);

-- name: GetComment :one
SELECT *
FROM comments
WHERE id = @id;

-- name: UpdateComment :exec
UPDATE comments
SET message = @message
WHERE id = @id;

-- name: DeleteComment :exec
DELETE
FROM comments
WHERE id = @id;

-- name: ListComments :many
SELECT *
FROM comments
WHERE ticket = @ticket
ORDER BY created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: CreateFeature :execlastid
INSERT INTO features (name)
VALUES (@name);

-- name: GetFeature :one
SELECT *
FROM features
WHERE id = @id;

-- name: UpdateFeature :exec
UPDATE features
SET name = @name
WHERE id = @id;

-- name: DeleteFeature :exec
DELETE
FROM features
WHERE id = @id;

-- name: ListFeatures :many
SELECT *
FROM features
ORDER BY created DESC;

------------------------------------------------------------------

-- name: CreateFile :execlastid
INSERT INTO files (name, blob, size, ticket)
VALUES (@name, @blob, @size, @ticket);

-- name: GetFile :one
SELECT *
FROM files
WHERE id = @id;

-- name: UpdateFile :exec
UPDATE files
SET name = @name, blob = @blob, size = @size
WHERE id = @id;

-- name: DeleteFile :exec
DELETE
FROM files
WHERE id = @id;

-- name: ListFiles :many
SELECT *
FROM files
WHERE ticket = @ticket
ORDER BY created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: CreateLink :execlastid
INSERT INTO links (name, url, ticket)
VALUES (@name, @url, @ticket);

-- name: GetLink :one
SELECT *
FROM links
WHERE id = @id;

-- name: UpdateLink :exec
UPDATE links
SET name = @name, url = @url
WHERE id = @id;

-- name: DeleteLink :exec
DELETE
FROM links
WHERE id = @id;

-- name: ListLinks :many
SELECT *
FROM links
WHERE ticket = @ticket
ORDER BY created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: CreateReaction :execlastid
INSERT INTO reactions (name, action, actiondata, trigger, triggerdata)
VALUES (@name, @action, @actiondata, @trigger, @triggerdata);

-- name: GetReaction :one
SELECT *
FROM reactions
WHERE id = @id;

-- name: UpdateReaction :exec
UPDATE reactions
SET name = @name, action = @action, actiondata = @actiondata, trigger = @trigger, triggerdata = @triggerdata
WHERE id = @id;

-- name: DeleteReaction :exec
DELETE
FROM reactions
WHERE id = @id;

-- name: ListReactions :many
SELECT *
FROM reactions
ORDER BY created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: CreateTask :execlastid
INSERT INTO tasks (name, open, owner, ticket)
VALUES (@name, @open, @owner, @ticket);

-- name: GetTask :one
SELECT *
FROM tasks
WHERE id = @id;

-- name: UpdateTask :exec
UPDATE tasks
SET name = @name, open = @open, owner = @owner
WHERE id = @id;

-- name: DeleteTask :exec
DELETE
FROM tasks
WHERE id = @id;

-- name: ListTasks :many
SELECT *
FROM tasks
WHERE ticket = @ticket
ORDER BY created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: CreateTimeline :execlastid
INSERT INTO timeline (message, ticket, time)
VALUES (@message, @ticket, @time);

-- name: GetTimeline :one
SELECT *
FROM timeline
WHERE id = @id;

-- name: UpdateTimeline :exec
UPDATE timeline
SET message = @message, time = @time
WHERE id = @id;

-- name: DeleteTimeline :exec
DELETE
FROM timeline
WHERE id = @id;

-- name: ListTimeline :many
SELECT *
FROM timeline
WHERE ticket = @ticket
ORDER BY created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: CreateType :execlastid
INSERT INTO types (singular, plural, icon, schema)
VALUES (@singular, @plural, @icon, @schema);

-- name: GetType :one
SELECT *
FROM types
WHERE id = @id;

-- name: UpdateType :exec
UPDATE types
SET singular = @singular, plural = @plural, icon = @icon, schema = @schema
WHERE id = @id;

-- name: DeleteType :exec
DELETE
FROM types
WHERE id = @id;

-- name: ListTypes :many
SELECT *
FROM types
ORDER BY created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: CreateUser :execlastid
INSERT INTO users (name, email, username, passwordHash, tokenKey)
VALUES (@name, @email, @username, @passwordHash, @tokenKey);

-- name: GetUser :one
SELECT *
FROM users
WHERE id = @id;

-- name: UpdateUser :exec
UPDATE users
SET name = @name, email = @email, username = @username, passwordHash = @passwordHash, tokenKey = @tokenKey
WHERE id = @id;

-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = @id;

-- name: ListUsers :many
SELECT *
FROM users
ORDER BY created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: CreateWebhook :execlastid
INSERT INTO webhooks (name, collection, destination)
VALUES (@name, @collection, @destination);

-- name: GetWebhook :one
SELECT *
FROM webhooks
WHERE id = @id;

-- name: UpdateWebhook :exec
UPDATE webhooks
SET name = @name, collection = @collection, destination = @destination
WHERE id = @id;

-- name: DeleteWebhook :exec
DELETE
FROM webhooks
WHERE id = @id;

-- name: ListWebhooks :many
SELECT *
FROM webhooks
ORDER BY created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: GetDashboardCounts :many
SELECT *
FROM dashboard_counts;

-- name: GetSidebar :many
SELECT *
FROM sidebar;

-- name: SearchTickets :many
SELECT *
FROM ticket_search
WHERE name LIKE '%' || @query || '%'
   OR description LIKE '%' || @query || '%'
   OR comment_messages LIKE '%' || @query || '%'
   OR file_names LIKE '%' || @query || '%'
   OR link_names LIKE '%' || @query || '%'
   OR link_urls LIKE '%' || @query || '%'
   OR task_names LIKE '%' || @query || '%'
   OR timeline_messages LIKE '%' || @query || '%'
LIMIT @limit OFFSET @offset;