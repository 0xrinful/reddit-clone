CREATE TABLE IF NOT EXISTS community_bans (
  community_id BIGINT NOT NULL REFERENCES communities (id) ON DELETE CASCADE,
  user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  banned_by BIGINT REFERENCES users (id) ON DELETE SET NULL,
  reason TEXT NOT NULL,
  created_at TIMESTAMP(0) WITH TIME ZONE DEFAULT NOW() NOT NULL,
  expires_at TIMESTAMP(0) WITH TIME ZONE,
  PRIMARY KEY (community_id, user_id)
);
