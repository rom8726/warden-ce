-- Create teams table
CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create team_members table
CREATE TABLE IF NOT EXISTS team_members (
    team_id INT NOT NULL REFERENCES teams(id),
    user_id INT NOT NULL REFERENCES users(id),
    role VARCHAR(50) NOT NULL, -- owner, admin, member
    PRIMARY KEY (team_id, user_id)
);

-- Create indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_teams_name ON teams(name);
CREATE INDEX IF NOT EXISTS idx_team_members_user_id ON team_members(user_id);

-- Add team_id column to projects table
ALTER TABLE projects ADD COLUMN team_id INT REFERENCES teams(id);
CREATE INDEX IF NOT EXISTS idx_projects_team_id ON projects(team_id);