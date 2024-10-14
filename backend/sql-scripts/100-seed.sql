\c bee;

SET search_path TO bee_schema, public;

-- Example values for E2E testing

-- user charlie

INSERT INTO bee_schema.users (id, username, access_token, refresh_token)
VALUES (-100, 'bee-ci-system', 'access_token', 'refresh_token');

INSERT INTO bee_schema.repos (id, name, user_id, latest_commit_sha, latest_commit_pushed_at, description)
VALUES (-200, 'example-using-beeci', -100, 'abc123def456ghi789', '2005-04-02 21:37:00', 'random desc 1'),
       (-201, 'example', -100, 'abc123def456ghi789', '2005-04-02 21:37:00', 'example repo!!! description');

INSERT INTO bee_schema.builds (repo_id, commit_sha, commit_message, installation_id, status)
VALUES (-200, '5ac1545229da7f0af6dfbae68950198866186a07', 'c_alpha commit 1', 0, 'queued');
INSERT INTO bee_schema.builds (repo_id, commit_sha, commit_message, installation_id, status)
VALUES (-200, '5ac1545229da7f0af6dfbae68950198866186a07', 'c_alpha commit 2', 0, 'queued');

INSERT INTO bee_schema.builds (repo_id, commit_sha, commit_message, installation_id, status)
VALUES (-201, '0262a10fb0590f29471feed5ecf53b418b5b0d67', 'c_bravo commit 1', 0, 'queued');

INSERT INTO bee_schema.builds (repo_id, commit_sha, commit_message, installation_id, status)
VALUES (-201, '0262a10fb0590f29471feed5ecf53b418b5b0d67', 'c_bravo commit 2', 0, 'queued');

-- user johnny

INSERT INTO bee_schema.users (id, username, access_token, refresh_token)
VALUES (-101, 'johnny', 'access_token', 'refresh_token');

INSERT INTO bee_schema.repos (id, name, user_id, latest_commit_sha, latest_commit_pushed_at, description)
VALUES (-203, 'j_alpha_repo', -101, 'abc123def456ghi789', '2005-04-02 21:37:00', 'random desc 1'),
       (-204, 'j_bravo_repo', -101, 'xyz987def456ghi789', '2005-04-02 21:37:00', 'random desc 1');

-- Insert a single build for each repo
INSERT INTO bee_schema.builds (repo_id, commit_sha, commit_message, installation_id, status)
VALUES (-203, '1234567890jkl', 'j_alpha commit 1', 1, 'queued'),
       (-204, '1234567890jkl', 'j_bravo commit 1', 1, 'queued');
