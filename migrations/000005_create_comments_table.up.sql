-- 1. Enable ltree extension
CREATE EXTENSION IF NOT EXISTS ltree;

-- 2. Define comment status enum
CREATE TYPE comment_status AS ENUM('visible', 'deleted', 'removed');

-- 3. Create the comments table
CREATE TABLE IF NOT EXISTS comments (
  id BIGSERIAL PRIMARY KEY,
  post_id BIGINT NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
  user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  parent_id BIGINT REFERENCES comments (id) ON DELETE CASCADE,
  content TEXT NOT NULL,
  score BIGINT NOT NULL DEFAULT 1,
  reply_count INT NOT NULL DEFAULT 0,
  depth INT NOT NULL DEFAULT 0,
  path LTREE,
  status comment_status NOT NULL DEFAULT 'visible',
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 4. Create performance indexes
CREATE INDEX idx_comments_path_gist ON comments USING GIST (path);
CREATE INDEX idx_comments_post_id ON comments (post_id);
CREATE INDEX idx_comments_path_btree ON comments USING BTREE (path);
CREATE INDEX idx_comments_post_depth_score ON comments (post_id, depth, score DESC);
CREATE INDEX idx_comments_parent_id ON comments (parent_id) WHERE parent_id IS NOT NULL;

-- 5. Define the Trigger Function
CREATE OR REPLACE FUNCTION handle_comment_hierarchy()
RETURNS TRIGGER AS $$
DECLARE
  v_parent_path LTREE;
  v_parent_depth INT;
BEGIN
  IF NEW.parent_id IS NULL THEN
      NEW.path := NEW.id::text::ltree;
      NEW.depth := 0;
  ELSE
      SELECT path, depth INTO v_parent_path, v_parent_depth
      FROM comments
      WHERE id = NEW.parent_id;

      -- Guard: safety check if parent record is missing
      IF v_parent_path IS NULL THEN
          RAISE EXCEPTION 'Parent comment with ID % does not exist', NEW.parent_id;
      END IF;

      -- Set the child's path: append child ID to parent path
      NEW.path := v_parent_path || NEW.id::text::ltree;
      NEW.depth := v_parent_depth + 1;

      -- update the parent's reply_count
      UPDATE comments 
      SET reply_count = reply_count + 1 
      WHERE id = NEW.parent_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 6. Attach the Trigger to the Table
CREATE TRIGGER comments_hierarchy_trigger
BEFORE INSERT ON comments
FOR EACH ROW
EXECUTE FUNCTION handle_comment_hierarchy();
