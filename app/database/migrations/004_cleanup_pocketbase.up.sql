DROP TABLE _migrations;
DROP TABLE _collections;
DROP TABLE _externalauths;

INSERT INTO users (id, avatar, created, email, emailVisibility, lastLoginAlertSentAt, lastResetSentAt,
                   lastVerificationSentAt, name, passwordHash, tokenKey, updated, username, verified)
SELECT id,
       avatar,
       created,
       email,
       false,
       '',
       lastResetSentAt,
       '',
       id,
       passwordHash,
       tokenKey,
       updated,
       id,
       true
FROM _admins;


INSERT INTO user_groups (user_id, group_id)
SELECT id, 'admin'
FROM _admins;

DROP TABLE _admins;

ALTER TABLE users
    DROP COLUMN emailvisibility;
