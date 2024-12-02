CREATE DATABASE "DisneyPlusChooser"
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;

CREATE SCHEMA IF NOT EXISTS disneyschema
    AUTHORIZATION postgres;

CREATE TABLE IF NOT EXISTS disneyschema.movies
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    title text COLLATE pg_catalog."default" NOT NULL,
    overview text COLLATE pg_catalog."default",
    horizontal_poster_w1080 text COLLATE pg_catalog."default",
    vertical_poster_w720 text COLLATE pg_catalog."default",
    CONSTRAINT movies_pkey PRIMARY KEY (id),
    CONSTRAINT unique_title UNIQUE (title)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS disneyschema.movies
    OWNER to postgres;

CREATE TABLE IF NOT EXISTS disneyschema.cache_runs
(
    id bigint NOT NULL,
    "timestamp" bigint NOT NULL,
    CONSTRAINT cache_runs_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS disneyschema.cache_runs
    OWNER to postgres;