BEGIN;

CREATE TABLE IF NOT EXISTS user_service.user_info
(
    user_id              UUID NOT NULL PRIMARY KEY,
    mail                 varchar                             NOT NULL
);

COMMIT;