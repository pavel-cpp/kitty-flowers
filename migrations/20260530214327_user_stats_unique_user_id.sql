-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_stats
    ADD CONSTRAINT user_stats_user_id_fk_unique
        UNIQUE (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_stats
    DROP CONSTRAINT user_stats_user_id_fk_unique;
-- +goose StatementEnd
