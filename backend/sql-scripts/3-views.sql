\c bee;

SET search_path TO bee_schema, public;

CREATE VIEW repos_and_users AS
SELECT
    repos.id AS repo_id,
    repos.name AS repo_name,
    repos.user_id,
    users.username
FROM bee_schema.repos
JOIN bee_schema.users ON users.id = repos.user_id;
