BEGIN;

DO
$$
    BEGIN
        PERFORM 'user_service.order_statuses'::regtype;
    EXCEPTION
        WHEN undefined_object THEN
            CREATE TYPE user_service.order_statuses AS ENUM (
                'created',
                'awaiting_payment',
                'payment_received',
                'completed',
                'canceled',
                'failed'
                );
    END
$$;

CREATE TABLE IF NOT EXISTS user_service.orders
(
    id                      serial PRIMARY KEY,
    user_id                 UUID                            NOT NULL,
    total_price             float                           NOT NULL,
    created_at              timestamp(0) with time zone     NOT NULL DEFAULT now(),
    status                  user_service.order_statuses     NOT NULL DEFAULT 'created'
);

COMMIT;
