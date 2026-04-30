DELETE FROM bee_schema.builds
WHERE repo_id IN (-200, -201, -202, -203, -204);

DELETE FROM bee_schema.repos
WHERE id IN (-200, -201, -202, -203, -204);

DELETE FROM bee_schema.users
WHERE id IN (-100, -101);
