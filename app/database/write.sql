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

-- name: CreateTicket :one
INSERT INTO tickets (id, name, description, open, owner, resolution, schema, state, type)
VALUES (@id, @name, @description, @open, @owner, @resolution, @schema, @state, @type)
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

-- name: CreateComment :one
INSERT INTO comments (id, author, message, ticket)
VALUES (@id, @author, @message, @ticket)
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

-- name: CreateFile :one
INSERT INTO files (id, name, blob, size, ticket)
VALUES (@id, @name, @blob, @size, @ticket)
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

-- name: CreateLink :one
INSERT INTO links (id, name, url, ticket)
VALUES (@id, @name, @url, @ticket)
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

-- name: CreateReaction :one
INSERT INTO reactions (id, name, action, actiondata, trigger, triggerdata)
VALUES (@id, @name, @action, @actiondata, @trigger, @triggerdata)
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

-- name: CreateTask :one
INSERT INTO tasks (id, name, open, owner, ticket)
VALUES (@id, @name, @open, @owner, @ticket)
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

-- name: CreateTimeline :one
INSERT INTO timeline (id, message, ticket, time)
VALUES (@id, @message, @ticket, @time)
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

-- name: CreateType :one
INSERT INTO types (id, singular, plural, icon, schema)
VALUES (@id, @singular, @plural, @icon, @schema)
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

-- name: CreateUser :one
INSERT INTO users (id, name, email, username, passwordHash, tokenKey, avatar, verified)
VALUES (@id, @name, @email, @username, @passwordHash, @tokenKey, @avatar, @verified)
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
    lastLoginAlertSentAt   = coalesce(sqlc.narg('lastLoginAlertSentAt'), lastLoginAlertSentAt),
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

-- name: CreateWebhook :one
INSERT INTO webhooks (id, name, collection, destination)
VALUES (@id, @name, @collection, @destination)
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

-- name: CreateGroup :one
INSERT INTO groups (id, name, permissions)
VALUES (@id, @name, @permissions)
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
