CREATE DATABASE bee;

CREATE TABLE users (
    id SERIAL PRIMARY KEY, -- GitHub user id
    username VARCHAR(255) UNIQUE NOT NULL, -- GitHub username
    installation_token VARCHAR(40) NOT NULL
);

CREATE TABLE builds (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    repo VARCHAR(255) NOT NULL,
    branch VARCHAR(255) NOT NULL,
    commit VARCHAR(40) NOT NULL,
    status VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
