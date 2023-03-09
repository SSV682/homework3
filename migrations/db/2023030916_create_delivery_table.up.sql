BEGIN;

CREATE TABLE IF NOT EXISTS user_service.delivery
(
    id                      serial PRIMARY KEY,
    order_id                int                             NOT NULL,
    date                    timestamp(0) with time zone     NOT NULL,
    order_content           jsonb                           NOT NULL,
    address                 jsonb                           NOT NULL
);

COMMIT;