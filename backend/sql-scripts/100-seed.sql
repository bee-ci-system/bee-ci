\c bee;

SET search_path TO bee_schema, public;

-- Example values for E2E testing

-- user charlie

INSERT INTO bee_schema.users (id, username, access_token, refresh_token)
VALUES (-100, 'charlie', 'access_token', 'refresh_token');

INSERT INTO bee_schema.repos (id, name, user_id)
VALUES (-200, 'c_alpha', -100),
       (-201, 'c_bravo', -100),
       (-202, 'c_charlie', -100);

INSERT INTO bee_schema.builds (repo_id, commit_sha, commit_message, installation_id, status)
VALUES (-200, '1234567890abc', 'c_alpha commit 1', 0, 'queued');
INSERT INTO bee_schema.builds (repo_id, commit_sha, commit_message, installation_id, status)
VALUES (-200, '1234567890xyz', 'c_alpha commit 2', 0, 'queued');

INSERT INTO bee_schema.builds (repo_id, commit_sha, commit_message, installation_id, status)
VALUES (-201, '1234567890def', 'c_bravo commit 1', 0, 'queued');

INSERT INTO bee_schema.builds (repo_id, commit_sha, commit_message, installation_id, status)
VALUES (-201, '1234567890ghi', 'c_bravo commit 2', 0, 'queued');

-- user johnny

INSERT INTO bee_schema.users (id, username, access_token, refresh_token)
VALUES (-101, 'johnny', 'access_token', 'refresh_token');

INSERT INTO bee_schema.repos (id, name, user_id)
VALUES (-203, 'j_alpha', -101);

INSERT INTO bee_schema.builds (repo_id, commit_sha, commit_message, installation_id, status)
VALUES (-203, '1234567890jkl', 'j_alpha commit 1', 1, 'queued');
