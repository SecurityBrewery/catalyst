CREATE TABLE roles
(
    id          TEXT default ('r' || lower(hex(randomblob(7)))) not null
        primary key,
    name        TEXT UNIQUE                                     NOT NULL,
    permissions TEXT                                            NOT NULL, -- JSON array string like '["read:article","write:article"]'
    created     TEXT default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null,
    updated     TEXT default (strftime('%Y-%m-%d %H:%M:%fZ'))   not null
);

CREATE TABLE user_roles
(
    user_id TEXT NOT NULL,
    role_id TEXT NOT NULL,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE
);

CREATE TABLE role_inheritance
(
    parent_role_id TEXT NOT NULL,
    child_role_id  TEXT NOT NULL,
    PRIMARY KEY (parent_role_id, child_role_id),
    FOREIGN KEY (parent_role_id) REFERENCES roles (id) ON DELETE CASCADE,
    FOREIGN KEY (child_role_id) REFERENCES roles (id) ON DELETE CASCADE
);

CREATE VIEW user_effective_roles AS
WITH RECURSIVE all_roles(user_id, role_id) AS (
    -- Direct roles
    SELECT ur.user_id, ur.role_id, 'direct' AS role_type
    FROM user_roles ur

    UNION

    -- Inherited roles
    SELECT ar.user_id, ri.parent_role_id, 'inherited' AS role_type
    FROM all_roles ar
             JOIN role_inheritance ri ON ri.child_role_id = ar.role_id)
SELECT user_id,
       role_id,
       json_group_array(role_type) AS role_types
FROM all_roles;

CREATE VIEW user_effective_permissions AS
SELECT uer.user_id,
       CAST(json_each.value AS TEXT) AS permission
FROM user_effective_roles uer
         JOIN roles r ON r.id = uer.role_id, json_each(r.permissions);
