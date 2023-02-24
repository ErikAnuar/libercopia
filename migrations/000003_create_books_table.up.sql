CREATE TABLE IF NOT EXISTS books
(
    id          bigserial primary key,
    created_at  timestamp(0) with time zone NOT NULL DEFAULT now(),
    title       text                        NOT NULL,
    year        text                        NOT NULL,
    author      text                        NOT NULL,
    genres      text[]                      NOT NULL,
    price       double precision            NOT NULL,
    description text                        NOT NULL,
    version     integer                     NOT NULL DEFAULT 1
);