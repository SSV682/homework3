BEGIN;

CREATE SCHEMA IF NOT EXISTS user_service;

CREATE TABLE IF NOT EXISTS user_service.users
(
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    username varchar(255) UNIQUE NOT NULL,
    firstname varchar NOT NULL,
    lastname varchar NOT NULL,
    email varchar NOT NULL,
    phone varchar NOT NULL,
    "password" varchar NOT NULL
);

COMMIT;