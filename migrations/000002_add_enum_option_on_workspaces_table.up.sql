ALTER TABLE workspaces
DROP CONSTRAINT IF EXISTS workspaces_privacy_check;

ALTER TABLE workspaces
ADD CONSTRAINT workspaces_privacy_check
CHECK (privacy IN ('private', 'public', 'team'));