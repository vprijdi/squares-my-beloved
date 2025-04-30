ALTER TABLE tasks 
ADD COLUMN completed BOOLEAN NOT NULL DEFAULT FALSE;

-- 2. Trigger function to auto-update 'updated_at'
CREATE OR REPLACE FUNCTION update_tasks_modified_time()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 3. Create trigger (runs on INSERT/UPDATE)
CREATE TRIGGER tasks_update_modified_time
BEFORE INSERT OR UPDATE ON tasks
FOR EACH ROW
EXECUTE FUNCTION update_tasks_modified_time();