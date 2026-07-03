CREATE TABLE IF NOT EXISTS community_members (
  community_id BIGINT NOT NULL REFERENCES communities (id) ON DELETE CASCADE,
  user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  role TEXT NOT NULL DEFAULT 'member',
  joined_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  PRIMARY KEY (community_id, user_id)
);
