CREATE TABLE IF NOT EXISTS communities (
  id BIGSERIAL PRIMARY KEY,
  name CITEXT NOT NULL UNIQUE,
  owner_id BIGINT REFERENCES users (id) ON DELETE SET NULL,
  description TEXT NOT NULL,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_communities_name_trgm ON communities USING gin (name gin_trgm_ops);
