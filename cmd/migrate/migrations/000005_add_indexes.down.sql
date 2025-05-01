BEGIN;

DROP INDEX IF EXISTS idx_tasks_user_id;
DROP INDEX IF EXISTS idx_tasks_created_at;
DROP INDEX IF EXISTS idx_tasks_user_recent;
DROP INDEX IF EXISTS idx_tasks_user_incomplete;

COMMIT;