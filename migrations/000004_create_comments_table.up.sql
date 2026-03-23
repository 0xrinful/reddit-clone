CREATE TABLE IF NOT EXISTS comments (
  id BIGSERIAL PRIMARY KEY,
  content TEXT NOT NULL,
  post_id BIGINT NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
  user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  parent_id BIGINT REFERENCES comments (id) ON DELETE CASCADE,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments (post_id);

CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments (parent_id);
