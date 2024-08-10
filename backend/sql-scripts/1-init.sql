CREATE DATABASE bee;
GRANT ALL PRIVILEGES ON DATABASE bee TO "postgres";

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

-- Taken from https://docs.github.com/en/rest/checks/runs?apiVersion=2022-11-28#create-a-check-run
CREATE TYPE build_status AS ENUM ('queued', 'in_progress', 'completed');
CREATE TYPE build_conclusion AS ENUM ('canceled', 'failure', 'success', 'timed_out');

CREATE TABLE bee_schema.builds (
    id SERIAL PRIMARY KEY, -- aka external_id for GitHub check run
    repo_id INTEGER NOT NULL,
    commit_sha VARCHAR(40) NOT NULL,
    status build_status NOT NULL,
    conclusion build_conclusion,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT status_completed_requires_conclusion CHECK (
        conclusion IS NULL OR status = 'completed'
    )
);

INSERT INTO bee_schema.users (id, access_token, refresh_token) VALUES (
    2137, 'access_token', 'refresh_token'
);

INSERT INTO bee_schema.builds (repo_id, commit_sha, status) VALUES (
    1, '1234567890abcdef', 'queued'
);

CREATE OR REPLACE FUNCTION BUILDS_TRIGGER() RETURNS TRIGGER AS
$$
    BEGIN
        PERFORM pg_notify('builds_channel', row_to_json(NEW)::TEXT);
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER builds_notify_trigger AFTER INSERT OR UPDATE ON bee_schema.builds
FOR EACH ROW EXECUTE FUNCTION BUILDS_TRIGGER();
