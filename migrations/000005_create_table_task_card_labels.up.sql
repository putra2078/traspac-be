CREATE TABLE task_card_labels (
    id SERIAL PRIMARY KEY,
    task_card_id INT NOT NULL,
    title VARCHAR(255) NULL,
    color VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_task_card_labels
    FOREIGN KEY (task_card_id)
    REFERENCES task_cards(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE
);