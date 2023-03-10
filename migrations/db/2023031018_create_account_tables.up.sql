BEGIN;

CREATE TABLE IF NOT EXISTS user_service.account
(
    id                    UUID                              NOT NULL PRIMARY KEY,
    amount                float                             NOT NULL
);

CREATE TABLE IF NOT EXISTS user_service.outbox
(
    id                   serial                           NOT NULL PRIMARY KEY,
    topic                varchar                          NOT NULL,
    message              jsonb                            NOT NULL
);

COMMIT;