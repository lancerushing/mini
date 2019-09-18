
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
------------------------------------------------------

CREATE USER mini_app WITH PASSWORD 'mini_app';

GRANT USAGE ON SCHEMA public TO mini_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT,UPDATE,DELETE ON TABLES TO mini_app;

------------------------------------------------------
-- DROP TABLE IF EXISTS public.users;
CREATE TABLE public.users
(
    uuid uuid NOT NULL,
    name text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,

    PRIMARY KEY (uuid),
    UNIQUE  (email)
);

