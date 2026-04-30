CREATE SCHEMA bee_schema;

CREATE TYPE bee_schema.build_status AS ENUM ('queued', 'in_progress', 'completed');
CREATE TYPE bee_schema.build_conclusion AS ENUM ('canceled', 'failure', 'success', 'timed_out');

CREATE TABLE bee_schema.users
(
    id       BIGINT PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE bee_schema.repos
(
    id      BIGINT PRIMARY KEY,
    name    VARCHAR(256) NOT NULL,
    user_id BIGINT       NOT NULL,
    FOREIGN KEY (user_id) REFERENCES bee_schema.users (id) ON DELETE CASCADE
);

CREATE TABLE bee_schema.builds
(
    id              SERIAL PRIMARY KEY,
    repo_id         BIGINT                      NOT NULL,
    commit_sha      VARCHAR(40)                 NOT NULL,
    commit_message  VARCHAR(2048)               NOT NULL,
    installation_id BIGINT                      NOT NULL,
    check_run_id    BIGINT,
    status          bee_schema.build_status     NOT NULL,
    conclusion      bee_schema.build_conclusion,
    created_at      TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (repo_id) REFERENCES bee_schema.repos (id) ON DELETE CASCADE,
    CONSTRAINT status_completed_requires_conclusion CHECK (
        conclusion IS NULL OR status = 'completed'
    )
);
