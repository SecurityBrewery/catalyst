DROP TABLE _migrations;
DROP TABLE _collections;

INSERT INTO users (avatar, created, email, emailVisibility, lastLoginAlertSentAt, lastResetSentAt,
                   lastVerificationSentAt, name, passwordHash, tokenKey, updated, username, verified)
SELECT avatar, created, email, false, '', lastResetSentAt,
       '', id, passwordHash, tokenKey, updated, id, true
FROM _admins;