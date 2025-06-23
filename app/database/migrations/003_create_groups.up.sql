CREATE TABLE groups
(
    id          TEXT default ('g' || lower(hex(randomblob(7)))) not null
        primary key,
    name        TEXT UNIQUE                                     NOT NULL,
    permissions TEXT                                            NOT NULL, -- JSON array string like '["read:article","write:article"]'
    created     TEXT default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null,
    updated     TEXT default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null
);

CREATE TABLE user_groups
(
    user_id  TEXT NOT NULL,
    group_id TEXT NOT NULL,
    PRIMARY KEY (user_id, group_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE
);

CREATE TABLE group_inheritance
(
    parent_group_id TEXT NOT NULL,
    child_group_id  TEXT NOT NULL,
    PRIMARY KEY (parent_group_id, child_group_id),
    FOREIGN KEY (parent_group_id) REFERENCES groups (id) ON DELETE CASCADE,
    FOREIGN KEY (child_group_id) REFERENCES groups (id) ON DELETE CASCADE
);

CREATE VIEW group_effective_groups AS
WITH RECURSIVE all_groups(child_group_id, parent_group_id, group_type)
                   AS (SELECT rr.child_group_id, rr.parent_group_id, 'direct' AS group_type
                       FROM group_inheritance rr
                       UNION
                       SELECT ar.child_group_id, ri.parent_group_id, 'indirect' AS group_type
                       FROM all_groups ar
                                JOIN group_inheritance ri ON ri.child_group_id = ar.parent_group_id)
SELECT child_group_id, parent_group_id, group_type
FROM all_groups;

CREATE VIEW group_effective_permissions AS
SELECT re.parent_group_id, CAST(json_each.value AS TEXT) AS permission
FROM group_effective_groups re
         JOIN groups r ON r.id = re.child_group_id, json_each(r.permissions);

CREATE VIEW user_effective_groups AS
WITH RECURSIVE all_groups(user_id, group_id, group_type) AS (
    -- Direct groups
    SELECT ur.user_id, ur.group_id, 'direct' AS group_type
    FROM user_groups ur

    UNION

    -- Inherited groups
    SELECT ar.user_id, ri.child_group_id, 'indirect' AS group_type
    FROM all_groups ar
             JOIN group_inheritance ri ON ri.parent_group_id = ar.group_id)
SELECT user_id,
       group_id,
       group_type
FROM all_groups;

CREATE VIEW user_effective_permissions AS
SELECT DISTINCT uer.user_id,
                CAST(json_each.value AS TEXT) AS permission
FROM user_effective_groups uer
         JOIN groups r ON r.id = uer.group_id, json_each(r.permissions);

INSERT INTO groups (id, name, permissions)
VALUES ('analyst', 'Analyst', '["type:read", "file:read", "ticket:read", "ticket:write", "user:read", "group:read"]'),
       ('admin', 'Admin', '["admin"]');

INSERT INTO user_groups (user_id, group_id)
SELECT id, 'analyst'
FROM users
WHERE id != 'system'; -- Exclude the system user

INSERT INTO user_groups (user_id, group_id)
VALUES ('system', 'admin'); -- Assign the admin group to the system user