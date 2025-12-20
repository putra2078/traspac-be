CREATE TABLE task_card_comments (
    id SERIAL PRIMARY KEY,
    task_card_id INT NOT NULL,
    comment TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_task_card_comments_task_card
        FOREIGN KEY (task_card_id)
        REFERENCES task_cards(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

CREATE INDEX idx_task_card_comments_task_card_id
    ON task_card_comments(task_card_id);
