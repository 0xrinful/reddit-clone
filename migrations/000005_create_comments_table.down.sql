DROP TRIGGER IF EXISTS comments_hierarchy_trigger ON comments;

DROP FUNCTION IF EXISTS handle_comment_hierarchy ();

DROP TABLE IF EXISTS comments;

DROP TYPE IF EXISTS comment_status;

DROP EXTENSION IF EXISTS ltree;
