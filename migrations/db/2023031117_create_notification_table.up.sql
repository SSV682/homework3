BEGIN;

CREATE TABLE IF NOT EXISTS user_service.notification
(
    id                      serial                              PRIMARY KEY,
    mail                    varchar                             NOT NULL,
    message                 varchar                             NOT NULL
);

COMMIT;