-- name: CreateTicket :execlastid
INSERT INTO tickets (name)
VALUES (@name);

-- name: Ticket :one
SELECT *
FROM tickets
WHERE id = @id;

-- name: UpdateTicket :exec
UPDATE tickets
SET name = @name
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