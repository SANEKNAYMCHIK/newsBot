CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    tg_chat_id BIGINT UNIQUE,
    tg_username VARCHAR(100),
    tg_first_name VARCHAR(100),
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255),
    role VARCHAR(20) DEFAULT 'user'
);