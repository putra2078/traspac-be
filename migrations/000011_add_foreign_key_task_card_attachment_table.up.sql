ALTER TABLE task_card_attachments ADD CONSTRAINT fk_task_card_attachments_task_card_id FOREIGN KEY (task_card_id) REFERENCES task_cards(id);
