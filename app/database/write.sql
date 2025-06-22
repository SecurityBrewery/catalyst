-- name: CreateParam :exec
INSERT INTO _params (id, key, value)
VALUES (@id, @key, @value)
RETURNING *;

-- name: UpdateParam :exec
UPDATE _params
SET value = @value
WHERE key = @key
RETURNING *;

------------------------------------------------------------------

-- name: InsertTicket :one
INSERT INTO tickets (id, name, description, open, owner, resolution, schema, state, type, created, updated)
VALUES (@id, @name, @description, @open, @owner, @resolution, @schema, @state, @type, @created, @updated)
RETURNING *;

-- name: CreateTicket :one
INSERT INTO tickets (name, description, open, owner, resolution, schema, state, type)
VALUES (@name, @description, @open, @owner, @resolution, @schema, @state, @type)
RETURNING *;

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

------------------------------------------------------------------

-- name: InsertComment :one
INSERT INTO comments (id, author, message, ticket, created, updated)
VALUES (@id, @author, @message, @ticket, @created, @updated)
RETURNING *;

-- name: CreateComment :one
INSERT INTO comments (author, message, ticket)
VALUES (@author, @message, @ticket)
RETURNING *;

-- name: UpdateComment :one
UPDATE comments
SET message = coalesce(sqlc.narg('message'), message)
WHERE id = @id
RETURNING *;

-- name: DeleteComment :exec
DELETE
FROM comments
WHERE id = @id;

------------------------------------------------------------------

-- name: CreateFeature :one
INSERT INTO features (name)
VALUES (@name)
RETURNING *;

-- name: UpdateFeature :one
UPDATE features
SET name = coalesce(sqlc.narg('name'), name)
WHERE id = @id
RETURNING *;

-- name: DeleteFeature :exec
DELETE
FROM features
WHERE id = @id;

------------------------------------------------------------------

-- name: InsertFile :one
INSERT INTO files (id, name, blob, size, ticket, created, updated)
VALUES (@id, @name, @blob, @size, @ticket, @created, @updated)
RETURNING *;

-- name: CreateFile :one
INSERT INTO files (name, blob, size, ticket)
VALUES (@name, @blob, @size, @ticket)
RETURNING *;

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

------------------------------------------------------------------

-- name: InsertLink :one
INSERT INTO links (id, name, url, ticket, created, updated)
VALUES (@id, @name, @url, @ticket, @created, @updated)
RETURNING *;

-- name: CreateLink :one
INSERT INTO links (name, url, ticket)
VALUES (@name, @url, @ticket)
RETURNING *;

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

------------------------------------------------------------------

-- name: InsertReaction :one
INSERT INTO reactions (id, name, action, actiondata, trigger, triggerdata, created, updated)
VALUES (@id, @name, @action, @actiondata, @trigger, @triggerdata, @created, @updated)
RETURNING *;

-- name: CreateReaction :one
INSERT INTO reactions (name, action, actiondata, trigger, triggerdata)
VALUES (@name, @action, @actiondata, @trigger, @triggerdata)
RETURNING *;

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

------------------------------------------------------------------

-- name: InsertTask :one
INSERT INTO tasks (id, name, open, owner, ticket, created, updated)
VALUES (@id, @name, @open, @owner, @ticket, @created, @updated)
RETURNING *;

-- name: CreateTask :one
INSERT INTO tasks (name, open, owner, ticket)
VALUES (@name, @open, @owner, @ticket)
RETURNING *;

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

------------------------------------------------------------------

-- name: InsertTimeline :one
INSERT INTO timeline (id, message, ticket, time, created, updated)
VALUES (@id, @message, @ticket, @time, @created, @updated)
RETURNING *;

-- name: CreateTimeline :one
INSERT INTO timeline (message, ticket, time)
VALUES (@message, @ticket, @time)
RETURNING *;

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

------------------------------------------------------------------

-- name: InsertType :one
INSERT INTO types (id, singular, plural, icon, schema, created, updated)
VALUES (@id, @singular, @plural, @icon, @schema, @created, @updated)
RETURNING *;

-- name: CreateType :one
INSERT INTO types (singular, plural, icon, schema)
VALUES (@singular, @plural, @icon, @schema)
RETURNING *;

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

------------------------------------------------------------------

-- name: InsertUser :one
INSERT INTO users (id, name, email, username, passwordHash, tokenKey, avatar, verified, created, updated)
VALUES (@id, @name, @email, @username, @passwordHash, @tokenKey, @avatar, @verified, @created, @updated)
RETURNING *;

-- name: CreateUser :one
INSERT INTO users (name, email, username, passwordHash, tokenKey, avatar, verified)
VALUES (@name, @email, @username, @passwordHash, @tokenKey, @avatar, @verified)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET name                   = coalesce(sqlc.narg('name'), name),
    email                  = coalesce(sqlc.narg('email'), email),
    username               = coalesce(sqlc.narg('username'), username),
    passwordHash           = coalesce(sqlc.narg('passwordHash'), passwordHash),
    tokenKey               = coalesce(sqlc.narg('tokenKey'), tokenKey),
    avatar                 = coalesce(sqlc.narg('avatar'), avatar),
    verified               = coalesce(sqlc.narg('verified'), verified),
    lastResetSentAt        = coalesce(sqlc.narg('lastResetSentAt'), lastResetSentAt),
    lastVerificationSentAt = coalesce(sqlc.narg('lastVerificationSentAt'), lastVerificationSentAt)
WHERE id = @id
  AND id != 'system'
RETURNING *;

-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = @id
  AND id != 'system';

------------------------------------------------------------------

-- name: InsertWebhook :one
INSERT INTO webhooks (id, name, collection, destination, created, updated)
VALUES (@id, @name, @collection, @destination, @created, @updated)
RETURNING *;

-- name: CreateWebhook :one
INSERT INTO webhooks (name, collection, destination)
VALUES (@name, @collection, @destination)
RETURNING *;

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

------------------------------------------------------------------

-- name: InsertGroup :one
INSERT INTO groups (id, name, permissions, created, updated)
VALUES (@id, @name, @permissions, @created, @updated)
RETURNING *;

-- name: CreateGroup :one
INSERT INTO groups (name, permissions)
VALUES (@name, @permissions)
RETURNING *;

-- name: UpdateGroup :one
UPDATE groups
SET name        = coalesce(sqlc.narg('name'), name),
    permissions = coalesce(sqlc.narg('permissions'), permissions)
WHERE id = @id
RETURNING *;

-- name: DeleteGroup :exec
DELETE
FROM groups
WHERE id = @id;

-- name: AssignGroupToUser :exec
INSERT INTO user_groups (user_id, group_id)
VALUES (@user_id, @group_id);

-- name: RemoveGroupFromUser :exec
DELETE
FROM user_groups
WHERE user_id = @user_id
  AND group_id = @group_id;

-- name: AssignParentGroup :exec
INSERT INTO group_inheritance (parent_group_id, child_group_id)
VALUES (@parent_group_id, @child_group_id);

-- name: RemoveParentGroup :exec
DELETE
FROM group_inheritance
WHERE parent_group_id = @parent_group_id
  AND child_group_id = @child_group_id;
