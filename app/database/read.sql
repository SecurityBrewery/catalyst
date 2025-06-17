-- name: Ticket :one
SELECT tickets.*, users.name as owner_name, types.singular as type_singular, types.plural as type_plural
FROM tickets
         LEFT JOIN users ON users.id = tickets.owner
         LEFT JOIN types ON types.id = tickets.type
WHERE tickets.id = @id;

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

------------------------------------------------------------------

-- name: GetComment :one
SELECT comments.*, users.name as author_name
FROM comments
         LEFT JOIN users ON users.id = comments.author
WHERE comments.id = @id;

-- name: ListComments :many
SELECT comments.*, users.name as author_name, COUNT(*) OVER () as total_count
FROM comments
         LEFT JOIN users ON users.id = comments.author
WHERE ticket = @ticket
   OR @ticket = ''
ORDER BY comments.created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: GetFeature :one
SELECT *
FROM features
WHERE id = @id;

-- name: ListFeatures :many
SELECT features.*, COUNT(*) OVER () as total_count
FROM features
ORDER BY features.created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: GetFile :one
SELECT *
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

-- name: GetLink :one
SELECT *
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

-- name: GetReaction :one
SELECT *
FROM reactions
WHERE id = @id;

-- name: ListReactions :many
SELECT reactions.*, COUNT(*) OVER () as total_count
FROM reactions
ORDER BY reactions.created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: GetTask :one
SELECT tasks.*, users.name as owner_name, tickets.name as ticket_name, tickets.type as ticket_type
FROM tasks
         LEFT JOIN users ON users.id = tasks.owner
         LEFT JOIN tickets ON tickets.id = tasks.ticket
WHERE tasks.id = @id;

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

-- name: GetTimeline :one
SELECT *
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

-- name: GetType :one
SELECT *
FROM types
WHERE id = @id;

-- name: ListTypes :many
SELECT types.*, COUNT(*) OVER () as total_count
FROM types
ORDER BY created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

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

-- name: ListUsers :many
SELECT users.*, COUNT(*) OVER () as total_count
FROM users
WHERE id != 'system'
ORDER BY users.created DESC
LIMIT @limit OFFSET @offset;

------------------------------------------------------------------

-- name: GetWebhook :one
SELECT *
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

-- name: GetGroup :one
SELECT *
FROM groups
WHERE id = @id;

-- name: ListGroups :many
SELECT g.*, COUNT(*) OVER () as total_count
FROM groups AS g
ORDER BY g.created DESC
LIMIT @limit OFFSET @offset;

-- name: ListUserGroups :many
SELECT g.*, uer.group_type, COUNT(*) OVER () as total_count
FROM user_effective_groups uer
         JOIN groups AS g ON g.id = uer.group_id
WHERE uer.user_id = @user_id
ORDER BY g.name DESC;

-- name: ListGroupUsers :many
SELECT users.*, uer.group_type
FROM user_effective_groups uer
         JOIN users ON users.id = uer.user_id
WHERE uer.group_id = @group_id
ORDER BY users.name DESC;

-- name: ListUserPermissions :many
SELECT user_effective_permissions.permission
FROM user_effective_permissions
WHERE user_id = @user_id
ORDER BY permission;

-- name: ListParentGroups :many
SELECT g.*, group_effective_groups.group_type
FROM group_effective_groups
         JOIN groups AS g ON g.id = group_effective_groups.child_group_id
WHERE parent_group_id = @group_id
ORDER BY group_effective_groups.group_type;

-- name: ListChildGroups :many
SELECT g.*, group_effective_groups.group_type
FROM group_effective_groups
         JOIN groups AS g ON g.id = group_effective_groups.parent_group_id
WHERE child_group_id = @group_id
ORDER BY group_effective_groups.group_type;

-- name: ListParentPermissions :many
SELECT group_effective_permissions.permission
FROM group_effective_permissions
WHERE parent_group_id = @group_id
ORDER BY permission;
