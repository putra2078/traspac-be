-- Alter boards table foreign key constraint from CASCADE to RESTRICT
ALTER TABLE boards
DROP CONSTRAINT IF EXISTS fk_workspace_boards;

ALTER TABLE boards
ADD CONSTRAINT fk_workspace_boards
FOREIGN KEY (workspace_id)
REFERENCES workspaces(id)
ON UPDATE CASCADE
ON DELETE RESTRICT;
