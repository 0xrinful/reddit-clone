CREATE TABLE IF NOT EXISTS refresh_tokens (
  id BIGSERIAL PRIMARY KEY,
  hash BYTEA NOT NULL UNIQUE,
  user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  scope TEXT NOT NULL,
  parent_id BIGINT REFERENCES refresh_tokens (id) ON DELETE SET NULL,
  family_id BIGINT NOT NULL REFERENCES refresh_tokens (id),
  used_at TIMESTAMP(0) WITH TIME ZONE,
  revoked_at TIMESTAMP(0) WITH TIME ZONE,
  expiry TIMESTAMP(0) WITH TIME ZONE NOT NULL,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id_scope ON refresh_tokens (user_id, scope);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_family_id ON refresh_tokens (family_id);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_scope ON refresh_tokens (user_id, scope);
