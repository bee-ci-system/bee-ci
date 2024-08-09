CREATE DATABASE bee;
GRANT ALL PRIVILEGES ON DATABASE bee to "postgres";

\c bee;

CREATE SCHEMA bee_schema;
SET search_path TO bee_schema, public;

CREATE TABLE bee_schema.users (
    id SERIAL PRIMARY KEY, -- GitHub user id
    username VARCHAR(255) UNIQUE NOT NULL, -- GitHub username
    installation_token VARCHAR(40) NOT NULL
);

-- CREATE TABLE bee_schema.repos (
--     id SERIAL PRIMARY KEY, -- GitHub repo id
--     user_id INTEGER NOT NULL,
--     name VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
--     FOREIGN KEY (user_id) REFERENCES bee_schema.users(id)
-- );

CREATE TABLE bee_schema.builds (
    id SERIAL PRIMARY KEY, -- aka external_id for GitHub check run
    repo_id VARCHAR(255) NOT NULL,
    commit_sha VARCHAR(40) NOT NULL, -- TODO: Use enum
    status VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- last updated_at is essentially completed_at
    -- FOREIGN KEY (repo_id) REFERENCES bee_schema.repos(id)
);

INSERT INTO bee_schema.users (username, installation_token) VALUES ('octocat', 'gho_1234567890');
INSERT INTO bee_schema.users (username, installation_token) VALUES ('octocat2', 'gho_1234567890sd');
