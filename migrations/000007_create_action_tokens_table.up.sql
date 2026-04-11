CREATE TABLE IF NOT EXISTS action_tokens (
  id BIGSERIAL PRIMARY KEY,
  hash BYTEA NOT NULL UNIQUE,
  user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  scope TEXT NOT NULL,
  expiry TIMESTAMP(0) WITH TIME ZONE NOT NULL,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_action_tokens_user_id_scope ON action_tokens (user_id, scope);
