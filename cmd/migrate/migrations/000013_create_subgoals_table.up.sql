CREATE TABLE subgoals (
    id BIGSERIAL PRIMARY KEY,
    goal_id BIGINT NOT NULL REFERENCES goals(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    position SMALLINT,                  
    achieved_at TIMESTAMPTZ,            
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_subgoals_goal_id ON subgoals(goal_id);
CREATE INDEX idx_subgoals_user_id ON subgoals(user_id);