BEGIN;

CREATE TABLE IF NOT EXISTS users
(
    id serial NOT NULL,
    username varchar(255),
    firstname varchar,
    lastname varchar,
    email varchar,
    phone varchar,
    password varchar,
    CONSTRAINT "users_pk" PRIMARY KEY (id)
    );

COMMIT;