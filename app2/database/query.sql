-- name: CreateTicket :one
INSERT INTO tickets (name, description, open, owner, resolution, schema, state, type)
VALUES (@name, @description, @open, @owner, @resolution, @schema, @state, @type)
RETURNING *;

-- name: Ticket :one
SELECT tickets.*, users.name as owner_name, types.singular as type_singular, types.plural as type_plural
FROM tickets
         LEFT JOIN users ON users.id = tickets.owner
         LEFT JOIN types ON types.id = tickets.type
WHERE tickets.id = @id;

-- name: UpdateTicket :one
UPDATE tickets
SET name        = coalesce(sqlc.narg('name'), name),
    description = coalesce(sqlc.narg('description'), description),
    open        = coalesce(sqlc.narg('open'), open),
    owner       = coalesce(sqlc.narg('owner'), owner),
    resolution  = coalesce(sqlc.narg('resolution'), resolution),
    schema      = coalesce(sqlc.narg('schema'), schema),
    state       = coalesce(sqlc.narg('state'), state),
    type        = coalesce(sqlc.narg('type'), type)
WHERE id = @id
RETURNING *;

-- name: DeleteTicket :exec
DELETE
FROM tickets
WHERE id = @id;

-- name: ListTickets :many
SELECT tickets.*, users.name as owner_name, types.singular as type_singular, types.plural as type_plural
FROM tickets
         LEFT JOIN users ON users.id = tickets.owner
         LEFT JOIN types ON types.id = tickets.type
ORDER BY tickets.created DESC
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

-- name: CreateComment :one
INSERT INTO comments (author, message, ticket)
VALUES (@author, @message, @ticket)
RETURNING *;

-- name: GetComment :one
SELECT comments.*, users.name as author_name
FROM comments
         LEFT JOIN users ON users.id = comments.author
WHERE comments.id = @id;

-- name: UpdateComment :one
UPDATE comments
SET message = coalesce(sqlc.narg('message'), message)
WHERE id = @id
RETURNING *;

-- name: DeleteComment :exec
DELETE
FROM comments
WHERE id = @id;

-- name: ListComments :many
SELECT comments.*, users.name as author_name
FROM comments
         LEFT JOIN users ON users.id = comments.author
WHERE ticket = @ticket
   OR @ticket = ''
ORDER BY comments.created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: CreateFeature :one
INSERT INTO features (name)
VALUES (@name)
RETURNING *;

-- name: GetFeature :one
SELECT *
FROM features
WHERE id = @id;

-- name: UpdateFeature :one
UPDATE features
SET name = coalesce(sqlc.narg('name'), name)
WHERE id = @id
RETURNING *;

-- name: DeleteFeature :exec
DELETE
FROM features
WHERE id = @id;

-- name: ListFeatures :many
SELECT *
FROM features
ORDER BY features.created DESC;

------------------------------------------------------------------

-- name: CreateFile :one
INSERT INTO files (name, blob, size, ticket)
VALUES (@name, @blob, @size, @ticket)
RETURNING *;

-- name: GetFile :one
SELECT *
FROM files
WHERE id = @id;

-- name: UpdateFile :one
UPDATE files
SET name = coalesce(sqlc.narg('name'), name),
    blob = coalesce(sqlc.narg('blob'), blob),
    size = coalesce(sqlc.narg('size'), size)
WHERE id = @id
RETURNING *;

-- name: DeleteFile :exec
DELETE
FROM files
WHERE id = @id;

-- name: ListFiles :many
SELECT *
FROM files
WHERE ticket = @ticket
   OR @ticket = ''
ORDER BY files.created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: CreateLink :one
INSERT INTO links (name, url, ticket)
VALUES (@name, @url, @ticket)
RETURNING *;

-- name: GetLink :one
SELECT *
FROM links
WHERE id = @id;

-- name: UpdateLink :one
UPDATE links
SET name = coalesce(sqlc.narg('name'), name),
    url  = coalesce(sqlc.narg('url'), url)
WHERE id = @id
RETURNING *;

-- name: DeleteLink :exec
DELETE
FROM links
WHERE id = @id;

-- name: ListLinks :many
SELECT *
FROM links
WHERE ticket = @ticket
   OR @ticket = ''
ORDER BY links.created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: CreateReaction :one
INSERT INTO reactions (name, action, actiondata, trigger, triggerdata)
VALUES (@name, @action, @actiondata, @trigger, @triggerdata)
RETURNING *;

-- name: GetReaction :one
SELECT *
FROM reactions
WHERE id = @id;

-- name: UpdateReaction :one
UPDATE reactions
SET name        = coalesce(sqlc.narg('name'), name),
    action      = coalesce(sqlc.narg('action'), action),
    actiondata  = coalesce(sqlc.narg('actiondata'), actiondata),
    trigger     = coalesce(sqlc.narg('trigger'), trigger),
    triggerdata = coalesce(sqlc.narg('triggerdata'), triggerdata)
WHERE id = @id
RETURNING *;

-- name: DeleteReaction :exec
DELETE
FROM reactions
WHERE id = @id;

-- name: ListReactions :many
SELECT *
FROM reactions
ORDER BY reactions.created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: CreateTask :one
INSERT INTO tasks (name, open, owner, ticket)
VALUES (@name, @open, @owner, @ticket)
RETURNING *;

-- name: GetTask :one
SELECT tasks.*, users.name as owner_name, tickets.name as ticket_name, tickets.type as ticket_type
FROM tasks
         LEFT JOIN users ON users.id = tasks.owner
         LEFT JOIN tickets ON tickets.id = tasks.ticket
WHERE tasks.id = @id;

-- name: UpdateTask :one
UPDATE tasks
SET name  = coalesce(sqlc.narg('name'), name),
    open  = coalesce(sqlc.narg('open'), open),
    owner = coalesce(sqlc.narg('owner'), owner)
WHERE id = @id
RETURNING *;

-- name: DeleteTask :exec
DELETE
FROM tasks
WHERE id = @id;

-- name: ListTasks :many
SELECT tasks.*, users.name as owner_name, tickets.name as ticket_name, tickets.type as ticket_type
FROM tasks
         LEFT JOIN users ON users.id = tasks.owner
         LEFT JOIN tickets ON tickets.id = tasks.ticket
WHERE ticket = @ticket
   OR @ticket = ''
ORDER BY tasks.created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: CreateTimeline :one
INSERT INTO timeline (message, ticket, time)
VALUES (@message, @ticket, @time)
RETURNING *;

-- name: GetTimeline :one
SELECT *
FROM timeline
WHERE id = @id;

-- name: UpdateTimeline :one
UPDATE timeline
SET message = coalesce(sqlc.narg('message'), message),
    time    = coalesce(sqlc.narg('time'), time)
WHERE id = @id
RETURNING *;

-- name: DeleteTimeline :exec
DELETE
FROM timeline
WHERE id = @id;

-- name: ListTimeline :many
SELECT *
FROM timeline
WHERE ticket = @ticket
   OR @ticket = ''
ORDER BY timeline.created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: CreateType :one
INSERT INTO types (singular, plural, icon, schema)
VALUES (@singular, @plural, @icon, @schema)
RETURNING *;

-- name: GetType :one
SELECT *
FROM types
WHERE id = @id;

-- name: UpdateType :one
UPDATE types
SET singular = coalesce(sqlc.narg('singular'), singular),
    plural   = coalesce(sqlc.narg('plural'), plural),
    icon     = coalesce(sqlc.narg('icon'), icon),
    schema   = coalesce(sqlc.narg('schema'), schema)
WHERE id = @id
RETURNING *;

-- name: DeleteType :exec
DELETE
FROM types
WHERE id = @id;

-- name: ListTypes :many
SELECT *
FROM types
ORDER BY created DESC;

------------------------------------------------------------------

-- name: CreateUser :one
INSERT INTO users (name, email, emailVisibility, username, passwordHash, tokenKey, avatar, verified)
VALUES (@name, @email, @emailVisibility, @username, @passwordHash, @tokenKey, @avatar, @verified)
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = @id;

-- name: UserByUserName :one
SELECT *
FROM users
WHERE username = @username;

-- name: UserByEmail :one
SELECT *
FROM users
WHERE email = @email;

-- name: UpdateUser :one
UPDATE users
SET name                   = coalesce(sqlc.narg('name'), name),
    email                  = coalesce(sqlc.narg('email'), email),
    emailVisibility        = coalesce(sqlc.narg('emailVisibility'), emailVisibility),
    username               = coalesce(sqlc.narg('username'), username),
    passwordHash           = coalesce(sqlc.narg('passwordHash'), passwordHash),
    tokenKey               = coalesce(sqlc.narg('tokenKey'), tokenKey),
    avatar                 = coalesce(sqlc.narg('avatar'), avatar),
    verified               = coalesce(sqlc.narg('verified'), verified)
WHERE id = @id
RETURNING *;

-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = @id;

-- name: ListUsers :many
SELECT *
FROM users
ORDER BY users.created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: CreateWebhook :one
INSERT INTO webhooks (name, collection, destination)
VALUES (@name, @collection, @destination)
RETURNING *;

-- name: GetWebhook :one
SELECT *
FROM webhooks
WHERE id = @id;

-- name: UpdateWebhook :one
UPDATE webhooks
SET name        = coalesce(sqlc.narg('name'), name),
    collection  = coalesce(sqlc.narg('collection'), collection),
    destination = coalesce(sqlc.narg('destination'), destination)
WHERE id = @id
RETURNING *;

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

-- name: ListSearchTickets :many
SELECT id,
       name,
       created,
       description,
       open,
       type,
       state,
       owner_name
FROM ticket_search
LIMIT @limit OFFSET @offset;

-- name: SearchTickets :many
SELECT id,
       name,
       created,
       description,
       open,
       type,
       state,
       owner_name
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