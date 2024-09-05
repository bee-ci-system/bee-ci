CREATE DATABASE bee;
GRANT ALL PRIVILEGES ON DATABASE bee TO "postgres";

\c bee;

CREATE SCHEMA bee_schema;
SET search_path TO bee_schema, public;

CREATE TABLE bee_schema.users
(
    id            INTEGER PRIMARY KEY,          -- GitHub user id
    username      VARCHAR(255) UNIQUE NOT NULL, -- username on GitHub (yes it can change, no we don't care)
    --installation_token VARCHAR(40) NOT NULL,
    access_token  VARCHAR(40)         NOT NULL,
    refresh_token VARCHAR(40)         NOT NULL
);

CREATE TABLE bee_schema.repos
(
    id      INTEGER PRIMARY KEY,   -- GitHub repo id
    name    VARCHAR(256) NOT NULL, -- name on GitHub (yes it can change, no we don't care)
    user_id INTEGER      NOT NULL,
    FOREIGN KEY (user_id) REFERENCES bee_schema.users (id) ON DELETE CASCADE
);

-- Taken from https://docs.github.com/en/rest/checks/runs?apiVersion=2022-11-28#create-a-check-run
CREATE TYPE build_status AS ENUM ('queued', 'in_progress', 'completed');
CREATE TYPE build_conclusion AS ENUM ('canceled', 'failure', 'success', 'timed_out');

CREATE TABLE bee_schema.builds
(
    id             SERIAL PRIMARY KEY,                -- aka external_id for GitHub check run
    repo_id        INTEGER                  NOT NULL,
    commit_sha     VARCHAR(40)              NOT NULL,
    commit_message VARCHAR(2048)            NOT NULl, -- we ain't handling longer commit messages
    status         build_status             NOT NULL,
    conclusion     build_conclusion,
    created_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (repo_id) REFERENCES bee_schema.repos (id) ON DELETE CASCADE,
    CONSTRAINT status_completed_requires_conclusion CHECK (
        conclusion IS NULL OR status = 'completed'
        )
);
