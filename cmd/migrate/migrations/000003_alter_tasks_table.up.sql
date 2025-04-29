ALTER TABLE tasks
ADD COLUMN updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
ADD COLUMN completions_count integer NOT NULL DEFAULT 0;