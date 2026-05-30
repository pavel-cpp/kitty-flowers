-- +goose Up
-- +goose StatementBegin
ALTER TABLE subscriptions
ALTER COLUMN next_run TYPE TIMESTAMPTZ;

ALTER TABLE subscriptions
ALTER COLUMN last_run TYPE TIMESTAMPTZ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE subscriptions
    ALTER COLUMN next_run TYPE TIMESTAMP;

ALTER TABLE subscriptions
    ALTER COLUMN last_run TYPE TIMESTAMP;
-- +goose StatementEnd
