CREATE TABLE IF NOT EXISTS users
(
    id serial NOT NULL,
    username varchar(255),
    firstname varchar,
    lastname varchar,
    email varchar,
    phone varchar,
    CONSTRAINT "users_pk" PRIMARY KEY (id)
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.users
    OWNER to postgres;