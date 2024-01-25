-- +goose Up
CREATE TABLE IF NOT EXISTS refresh_task
(
    id              SERIAL PRIMARY KEY,
    base            varchar(3),
    counter         varchar(3),
    ratio           DOUBLE PRECISION,
    time            timestamp,
    status          varchar NOT NULL DEFAULT 'in_progress',
    last_attempt_at timestamp
);

CREATE TABLE IF NOT EXISTS quote
(
    id      SERIAL PRIMARY KEY,
    base    varchar(3),
    counter varchar(3),
    ratio   DOUBLE PRECISION,
    time    timestamp,
    UNIQUE (base, counter)
);

-- +goose StatementBegin
-- +goose StatementEnd


-- +goose Down
DROP TABLE refresh_task;
DROP TABLE quote;