-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    last_run TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    next_run TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    active BOOLEAN DEFAULT TRUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS subscriptions;
-- +goose StatementEnd
