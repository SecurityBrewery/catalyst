-- name: CreateTicket :one
INSERT INTO tickets (id, name, description, open, owner, resolution, schema, state, type)
VALUES (@id, @name, @description, @open, @owner, @resolution, @schema, @state, @type)
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
SELECT tickets.*,
       users.name       as owner_name,
       types.singular   as type_singular,
       types.plural     as type_plural,
       COUNT(*) OVER () as total_count
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
INSERT INTO comments (id, author, message, ticket)
VALUES (@id, @author, @message, @ticket)
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
SELECT comments.*, users.name as author_name, COUNT(*) OVER () as total_count
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
SELECT features.*, COUNT(*) OVER () as total_count
FROM features
ORDER BY features.created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: CreateFile :one
INSERT INTO files (id, name, blob, size, ticket)
VALUES (@id, @name, @blob, @size, @ticket)
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
SELECT files.*, COUNT(*) OVER () as total_count
FROM files
WHERE ticket = @ticket
   OR @ticket = ''
ORDER BY files.created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: CreateLink :one
INSERT INTO links (id, name, url, ticket)
VALUES (@id, @name, @url, @ticket)
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
SELECT links.*, COUNT(*) OVER () as total_count
FROM links
WHERE ticket = @ticket
   OR @ticket = ''
ORDER BY links.created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: CreateReaction :one
INSERT INTO reactions (id, name, action, actiondata, trigger, triggerdata)
VALUES (@id, @name, @action, @actiondata, @trigger, @triggerdata)
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
SELECT reactions.*, COUNT(*) OVER () as total_count
FROM reactions
ORDER BY reactions.created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: CreateTask :one
INSERT INTO tasks (id, name, open, owner, ticket)
VALUES (@id, @name, @open, @owner, @ticket)
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
SELECT tasks.*,
       users.name       as owner_name,
       tickets.name     as ticket_name,
       tickets.type     as ticket_type,
       COUNT(*) OVER () as total_count
FROM tasks
         LEFT JOIN users ON users.id = tasks.owner
         LEFT JOIN tickets ON tickets.id = tasks.ticket
WHERE ticket = @ticket
   OR @ticket = ''
ORDER BY tasks.created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: CreateTimeline :one
INSERT INTO timeline (id, message, ticket, time)
VALUES (@id, @message, @ticket, @time)
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
SELECT timeline.*, COUNT(*) OVER () as total_count
FROM timeline
WHERE ticket = @ticket
   OR @ticket = ''
ORDER BY timeline.created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: CreateType :one
INSERT INTO types (id, singular, plural, icon, schema)
VALUES (@id, @singular, @plural, @icon, @schema)
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
SELECT types.*, COUNT(*) OVER () as total_count
FROM types
ORDER BY created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: CreateUser :one
INSERT INTO users (id, name, email, emailVisibility, username, passwordHash, tokenKey, avatar, verified)
VALUES (@id, @name, @email, @emailVisibility, @username, @passwordHash, @tokenKey, @avatar, @verified)
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = @id
  AND id != 'system';

-- name: UserByUserName :one
SELECT *
FROM users
WHERE username = @username
  AND id != 'system';

-- name: UserByEmail :one
SELECT *
FROM users
WHERE email = @email
  AND id != 'system';

-- name: SystemUser :one
SELECT *
FROM users
WHERE id = 'system';

-- name: UpdateUser :one
UPDATE users
SET name                   = coalesce(sqlc.narg('name'), name),
    email                  = coalesce(sqlc.narg('email'), email),
    emailVisibility        = coalesce(sqlc.narg('emailVisibility'), emailVisibility),
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

-- name: ListUsers :many
SELECT users.*, COUNT(*) OVER () as total_count
FROM users
WHERE id != 'system'
ORDER BY users.created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: CreateWebhook :one
INSERT INTO webhooks (id, name, collection, destination)
VALUES (@id, @name, @collection, @destination)
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
SELECT webhooks.*, COUNT(*) OVER () as total_count
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
SELECT id,
       name,
       created,
       description,
       open,
       type,
       state,
       owner_name,
       COUNT(*) OVER () as total_count
FROM ticket_search
WHERE (@query = '' OR (name LIKE '%' || @query || '%'
    OR description LIKE '%' || @query || '%'
    OR comment_messages LIKE '%' || @query || '%'
    OR file_names LIKE '%' || @query || '%'
    OR link_names LIKE '%' || @query || '%'
    OR link_urls LIKE '%' || @query || '%'
    OR task_names LIKE '%' || @query || '%'
    OR timeline_messages LIKE '%' || @query || '%'))
  AND (sqlc.narg('type') IS NULL OR type = sqlc.narg('type'))
  AND (sqlc.narg('open') IS NULL OR open = sqlc.narg('open'))
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: CreateRole :one
INSERT INTO roles (id, name, permissions)
VALUES (@id, @name, @permissions)
RETURNING *;

-- name: GetRole :one
SELECT roles.*, COUNT(*) OVER () as total_count
FROM roles
WHERE id = @id;

-- name: UpdateRole :one
UPDATE roles
SET name        = coalesce(sqlc.narg('name'), name),
    permissions = coalesce(sqlc.narg('permissions'), permissions)
WHERE id = @id
RETURNING *;

-- name: DeleteRole :exec
DELETE
FROM roles
WHERE id = @id;

-- name: ListRoles :many
SELECT roles.*, COUNT(*) OVER () as total_count
FROM roles
ORDER BY created DESC
LIMIT @limit OFFSET @offset;

-- name: AssignRoleToUser :exec
INSERT INTO user_roles (user_id, role_id)
VALUES (@user_id, @role_id);

-- name: RemoveRoleFromUser :exec
DELETE
FROM user_roles
WHERE user_id = @user_id
  AND role_id = @role_id;

-- name: AssignParentRole :exec
INSERT INTO role_inheritance (parent_role_id, child_role_id)
VALUES (@parent_role_id, @child_role_id);

-- name: RemoveParentRole :exec
DELETE
FROM role_inheritance
WHERE parent_role_id = @parent_role_id
  AND child_role_id = @child_role_id;

-- name: ListUserRoles :many
SELECT roles.*, uer.role_type, COUNT(*) OVER () as total_count
FROM user_effective_roles uer
         JOIN roles ON roles.id = uer.role_id
WHERE uer.user_id = @user_id
ORDER BY roles.name DESC;

-- name: ListUserPermissions :many
SELECT user_effective_permissions.permission
FROM user_effective_permissions
WHERE user_id = @user_id
ORDER BY permission;