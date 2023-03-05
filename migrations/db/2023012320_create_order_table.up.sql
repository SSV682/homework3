BEGIN;

CREATE TABLE IF NOT EXISTS user_service.orders
(
    id                      serial PRIMARY KEY,
    user_id                 UUID                            NOT NULL,
    total_price             float                           NOT NULL,
    products                jsonb                           NOT NULL,
    delivery_at             timestamp(0) with time zone     NOT NULL DEFAULT now(),
    created_at              timestamp(0) with time zone     NOT NULL DEFAULT now(),
    status                  varchar                         NOT NULL DEFAULT 'created'
);

COMMIT;
