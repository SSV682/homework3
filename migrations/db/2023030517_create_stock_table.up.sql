BEGIN;

CREATE TABLE IF NOT EXISTS user_service.products
(
    id                      serial PRIMARY KEY,
    quantity                   int                             NOT NULL,
    name                    varchar                         NOT NULL
);

COMMIT;