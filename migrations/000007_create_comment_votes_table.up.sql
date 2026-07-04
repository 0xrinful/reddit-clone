CREATE TABLE IF NOT EXISTS comment_votes (
  user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  comment_id BIGINT NOT NULL REFERENCES comments (id) ON DELETE CASCADE,
  value SMALLINT NOT NULL CHECK (value IN (-1, 1)),
  PRIMARY KEY (user_id, comment_id)
);
