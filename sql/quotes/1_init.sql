-- +goose Up
CREATE TABLE IF NOT EXISTS refresh_task
(
    id          SERIAL PRIMARY KEY,
    base        currency,
    counter     currency,
    ratio       DOUBLE PRECISION,
    time        timestamp,
    finished_at timestamp default null
);

CREATE TABLE IF NOT EXISTS quote
(
    id      SERIAL PRIMARY KEY,
    base    currency,
    counter currency,
    ratio       DOUBLE PRECISION,
    time    timestamp
);

-- +goose StatementBegin
-- +goose StatementEnd


-- +goose Down
DROP TABLE refresh_task;
DROP TABLE quote;
DROP TYPE currency;