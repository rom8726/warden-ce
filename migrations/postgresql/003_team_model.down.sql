-- Remove team_id column from projects table
DROP INDEX IF EXISTS idx_projects_team_id;
ALTER TABLE projects DROP COLUMN IF EXISTS team_id;

-- Drop indexes
DROP INDEX IF EXISTS idx_team_members_user_id;
DROP INDEX IF EXISTS idx_teams_name;

-- Drop team_members table
DROP TABLE IF EXISTS team_members;

-- Drop teams table
DROP TABLE IF EXISTS teams;