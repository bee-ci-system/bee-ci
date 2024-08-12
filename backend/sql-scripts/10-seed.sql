\c bee;

SET search_path TO bee_schema, public;

-- Example values for E2E testing

INSERT INTO bee_schema.users (id, username, access_token, refresh_token)
VALUES (-100, 'charlie', 'access_token', 'refresh_token');

INSERT INTO bee_schema.repos (id, name, user_id)
VALUES (-69, 'example_repo', -100);

INSERT INTO bee_schema.builds (repo_id, commit_sha, commit_message, status)
VALUES (-69, '1234567890abc', 'example commit', 'queued');

INSERT INTO bee_schema.builds (repo_id, commit_sha, commit_message, status)
VALUES (-69, '1234567890xyz', 'another example commit', 'queued');

