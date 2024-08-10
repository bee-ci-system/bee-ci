CREATE DATABASE bee;
GRANT ALL PRIVILEGES ON DATABASE bee to "postgres";

\c bee;

CREATE SCHEMA bee_schema;
SET search_path TO bee_schema, public;

CREATE TABLE bee_schema.users (
    id INTEGER PRIMARY KEY, -- GitHub user id
    -- username VARCHAR(255) UNIQUE NOT NULL, -- GitHub username
    --installation_token VARCHAR(40) NOT NULL,
    access_token VARCHAR(40) NOT NULL,
    refresh_token VARCHAR(40) NOT NULL
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

INSERT INTO bee_schema.users (id, access_token, refresh_token) VALUES (2137, 'access_token', 'refresh_token');

INSERT INTO bee_schema.builds (repo_id, commit_sha, status) VALUES ('octocat/hello-world', '1234567890abcdef', 'queued');

CREATE OR REPLACE FUNCTION builds_trigger() RETURNS TRIGGER AS
$$
    BEGIN
        PERFORM pg_notify('builds_channel', row_to_json(NEW)::TEXT);
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER builds_notify_trigger AFTER INSERT OR UPDATE ON bee_schema.builds
FOR EACH ROW EXECUTE FUNCTION builds_trigger();
