-- table users --
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE NULL
);

-- table contacts --
CREATE TABLE contacts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    photo VARCHAR(500) NOT NULL,
    email TEXT UNIQUE NOT NULL,
    phone_number VARCHAR(20),
    address TEXT,
    gender VARCHAR(15) CHECK (gender IN ('male', 'female', 'other')),
    birth_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


-- table settings --
CREATE TABLE settings (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    language VARCHAR(15) CHECK (language IN ('en_us', 'indo')),
    notification VARCHAR(15) CHECK (notification IN ('allowed', 'not_allowed')),

    CONSTRAINT fk_users_settings
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE
);


-- table workspaces --
CREATE TABLE workspaces (
    id SERIAL PRIMARY KEY,
    pass_code VARCHAR(255) NULL,
    created_by INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    privacy VARCHAR(15) CHECK (privacy IN ('private', 'public')),
    join_link VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_workspace_creator
    FOREIGN KEY (created_by)
    REFERENCES users(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE
);

CREATE TABLE boards (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	workspace_id INT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_workspace_boards
    FOREIGN KEY (workspace_id)
    REFERENCES workspaces(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
);

-- table task_tab --
CREATE TABLE task_tabs (
    id SERIAL PRIMARY KEY,
    board_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    position INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_boards_task_tab
    FOREIGN KEY (board_id)
    REFERENCES boards(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
);

-- table task_card --
CREATE TABLE task_cards (
    id SERIAL PRIMARY KEY,
    task_tab_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    content VARCHAR(255) NULL,
    comment VARCHAR(255) NULL,
    date DATE NOT NULL,
    status BOOLEAN NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_task_tab
    FOREIGN KEY (task_tab_id)
    REFERENCES task_tabs(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT

);

CREATE TABLE task_card_users (
    id SERIAL PRIMARY KEY,
    task_card_id INT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_task_card
    FOREIGN KEY (task_card_id)
    REFERENCES task_cards(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT,

    CONSTRAINT fk_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
);

-- table servers --
CREATE TABLE servers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_by INT NOT NULL,
    privacy VARCHAR(15) CHECK (privacy IN ('public', 'private')),
    pass_code VARCHAR(255) NULL,
    link_join VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_server_creator
    FOREIGN KEY (created_by)
    REFERENCES users(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
);

-- table room_chat --
CREATE TABLE rooms_chats (
    id SERIAL PRIMARY KEY,
    server_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_by INT NOT NULL,
    pass_code VARCHAR(255) NULL,
    link_join VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_rooms_creator
    FOREIGN KEY (created_by)
    REFERENCES users(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT,

    CONSTRAINT fk_room_server
    FOREIGN KEY (server_id)
    REFERENCES servers(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
);

-- table room_messages --
CREATE TABLE room_messages (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE RESTRICT,
    room_id INT NOT NULL,
    message_text VARCHAR(500) NOT NULL,
    message_content VARCHAR(500) NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_messages_room
    FOREIGN KEY(room_id)
    REFERENCES rooms_chats(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
);

-- table direct_message --
CREATE TABLE direct_messages (
    id SERIAL PRIMARY KEY,
    user_receiver INTEGER REFERENCES users(id) ON DELETE RESTRICT,
    user_sender INTEGER REFERENCES users(id) ON DELETE RESTRICT,
    message_text VARCHAR(500) NOT NULL,
    message_content VARCHAR(500) NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);