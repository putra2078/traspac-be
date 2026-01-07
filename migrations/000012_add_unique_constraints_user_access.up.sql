ALTER TABLE task_card_users ADD CONSTRAINT unique_task_card_user UNIQUE (task_card_id, user_id);
ALTER TABLE boards_users ADD CONSTRAINT unique_board_user UNIQUE (board_id, user_id);
ALTER TABLE workspaces_users ADD CONSTRAINT unique_workspace_user UNIQUE (workspace_id, user_id);
