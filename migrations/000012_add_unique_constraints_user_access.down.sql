ALTER TABLE task_card_users DROP CONSTRAINT IF EXISTS unique_task_card_user;
ALTER TABLE boards_users DROP CONSTRAINT IF EXISTS unique_board_user;
ALTER TABLE workspaces_users DROP CONSTRAINT IF EXISTS unique_workspace_user;

