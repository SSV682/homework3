BEGIN;

CREATE TABLE IF NOT EXISTS user_service.notification
(
    id                      serial                              PRIMARY KEY,
    user_id                 varchar                             NOT NULL,
    message                 varchar                             NOT NULL
);

COMMIT;