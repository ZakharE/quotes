-- +goose Up
CREATE TABLE IF NOT EXISTS quote
(
    id      SERIAL PRIMARY KEY,
    base    varchar(3),
    counter varchar(3),
    ratio   DOUBLE PRECISION,
    time    timestamp,
    UNIQUE (base, counter)
);

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

CREATE INDEX idx_refresh_task_last_attempt_status
    ON refresh_task (last_attempt_at, status);

-- +goose StatementBegin
-- +goose StatementEnd


-- +goose Down
DROP TABLE refresh_task;
DROP TABLE quote;