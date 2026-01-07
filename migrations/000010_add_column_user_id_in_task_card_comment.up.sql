ALTER TABLE task_card_comments ADD COLUMN user_id INTEGER;
ALTER TABLE task_card_comments ADD CONSTRAINT fk_task_card_comments_user_id FOREIGN KEY (user_id) REFERENCES users(id);