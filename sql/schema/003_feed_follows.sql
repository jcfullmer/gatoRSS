-- +goose Up
CREATE TABLE feed_follows(
 id UUID PRIMARY KEY,
 created_at TIMESTAMP NOT NULL,
 updated_at TIMESTAMP NOT NULL,
 user_id UUID NOT NULL CONSTRAINT fk_user_id REFERENCES users(id) ON DELETE CASCADE,
 feed_id UUID NOT NULL CONSTRAINT fk_feed_id REFERENCES feeds(id) ON DELETE CASCADE,
 CONSTRAINT uq_user_feed_id UNIQUE (user_id, feed_id)
 );
 
-- +goose Down
DROP TABLE feed_follows;