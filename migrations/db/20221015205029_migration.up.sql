BEGIN;

CREATE TABLE IF NOT EXISTS users
(
    id serial PRIMARY KEY NOT NULL,
    username varchar(255) UNIQUE NOT NULL,
    firstname varchar NOT NULL,
    lastname varchar NOT NULL,
    email varchar NOT NULL,
    phone varchar NOT NULL,
    "password" varchar NOT NULL
);

COMMIT;