BEGIN;

CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id);
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);
CREATE INDEX IF NOT EXISTS idx_tasks_user_recent ON tasks(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_tasks_user_incomplete ON tasks(user_id) 
WHERE completions_count > 0;

COMMIT;