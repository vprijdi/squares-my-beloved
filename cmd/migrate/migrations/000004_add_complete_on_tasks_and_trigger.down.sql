DROP TRIGGER IF EXISTS tasks_update_modified_time ON tasks;

DROP FUNCTION IF EXISTS update_tasks_modified_time();

ALTER TABLE tasks 
DROP COLUMN IF EXISTS completed;