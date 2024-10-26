# backend

We try to keep dependencies to a **reasonable** minimum.

[Project layout](https://go.dev/doc/modules/layout)

## Auth

[GitHub Apps have three methods of authentication][link_1]. We support:

- [x] [Authenticating as a GitHub App][auth_gh_app] using a JSON Web Token (JWT)
- [x] [Authenticating as a specific installation][auth_gh_install] of a GitHub
  App using an installation access token
- [ ] [Authenticating on behalf of a user][auth_gh_user]

### How does auth work?

We use a thing called [web application flow].
It's also explained from the practical side in [Build a "Login" button tutorial][login_btn].

Overall, it looks like this:

- The user clicks the "Sign in with GitHub" button and is redirected
  [https://github.com/login/oauth/authorize?client_id=Iv23liiZSvMGEpgOlexa]([https://github.com/login/oauth/authorize?client_id=Iv23liiZSvMGEpgOlexa])

  ![](./assets/demo.avif)

- Once user clicks "Authorize Bee CI system", they are redirected to our website

  `$BASE_URL/webhook/github/callback?code={code}`. Upon visiting that site, our
  backend get the `code` from URL params and sends POST to
  `https://github.com/login/oauth/access_token`, with the following payload:

  ```json
  {
    "client_id": "$CLIENT_ID",
    "client_secret": "$CLIENT_SECRET",
    "code": "${code}"
  }
  ```

- GitHub responds with an **access token**.

### What is that access token?

It looks like this:

```
ghu_Nr8ecJD8nWC4DKRvK694YpS8uJ5oHl0ix0sN
```

[TENTATIVE]

- It's something like a PAT (Personal Access Token).
- We persistently save it to a database and NEVER expose it to the user.
  It's an important and sensitive secret!
- Whenever we get a new token, we need to replace the old one.

See also:

- [GitHub's token format]

### What can the access token be used for?

The access token lets us make requests to the API on behalf of a user.

For example, we can use it to access user's private repositories:

```console
curl \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  https://api.github.com/repos/bartekpacia/discover_rudy
```

This token can be entirely unrelated to how sign in works in our webapp.

Our webapp can issue a JWT that is completely independent of the access token.

## Storage

### Postgres database

We don't use an ORM, instead we use [sqlx](https://jmoiron.github.io/sqlx).

[pgcli](https://www.pgcli.com) is a good client for PostgreSQL.

#### Connect to the local database

```console
pgcli -h localhost -p 5432 -u postgres -W -d bee
```

#### Connect to the production database

The fastest way to get the connection string:

```console
$ terraform output db_postgres_uri
"postgresql://doadmin:AVNS_VK_8D3KFy3dA2_8NjXw@bee-db-cluster-postgres-do-user-18862315-0.l.db.ondigitalocean.com:25060/defaultdb?sslmode=require"
```

Copy and paste it for `pgcli`, replacing `defaultdb` with `bee`:

```console
pgcli "postgresql://doadmin:AVNS_VK_8D3KFy3dA2_8NjXw@bee-db-cluster-postgres-do-user-18862315-0.l.db.ondigitalocean.com:25060/defaultdb?sslmode=require"
```

Then set the search path:

```postgresql
SET search_path TO bee_schema, public;
```

Now you're ready to interact with the database, for example:

```postgresql
SELECT builds.id AS build_id, builds.created_at, repo_id, repos.name AS repo_name, commit_sha, commit_message
FROM bee_schema.builds
         JOIN bee_schema.repos ON repos.id = builds.repo_id;
```

| build_id | created_at                    | repo_id   | repo_name           | commit_sha                               | commit_message                    |
|----------|-------------------------------|-----------|---------------------|------------------------------------------|-----------------------------------|
| 1        | 2024-10-22 17:30:21.268332+00 | -200      | example-using-beeci | 5ac1545229da7f0af6dfbae68950198866186a07 | c_alpha commit 1                  |
| 2        | 2024-10-22 17:30:21.421748+00 | -200      | example-using-beeci | 5ac1545229da7f0af6dfbae68950198866186a07 | c_alpha commit 2                  |
| 3        | 2024-10-22 17:30:21.571329+00 | -201      | example-using-beeci | 0262a10fb0590f29471feed5ecf53b418b5b0d67 | c_bravo commit 1                  |
| 4        | 2024-10-22 17:30:21.72085+00  | -201      | example-using-beeci | 0262a10fb0590f29471feed5ecf53b418b5b0d67 | c_bravo commit 2                  |
| 5        | 2024-10-22 17:30:22.176884+00 | -203      | j_alpha_repo        | 1234567890jkl                            | j_alpha commit 1                  |
| 6        | 2024-10-22 17:30:22.176884+00 | -204      | j_bravo_repo        | 1234567890jkl                            | j_bravo commit 1                  |
| 9        | 2024-10-24 16:01:24.351654+00 | 830117435 | bee-ci              | 6189e145baec92b840aeafd071264f7520afd1a2 | feat/added-error-page             |
| 10       | 2024-10-24 18:44:22.098719+00 | 830117435 | bee-ci              | d222438db174fa82db24cda37a179efad07ff589 | feat/added-sad-bee-for-error-page |

### Redis database

```console
$ terraform output db_redis_uri
"rediss://default:AVNS_REDACTEDvQkTk0HB44-@bee-db-cluster-redis-do-user-12345678-0.l.db.ondigitalocean.com:25061"
```

Copy and paste it for `redis-cli`:

```console
redis-cli -u "rediss://default:AVNS_REDACTEDvQkTk0HB44-@bee-db-cluster-redis-do-user-12345678-0.l.db.ondigitalocean.com:25061"
```

Now you're ready to interact with the database, for example:

```redis
KEYS "*"
```

### Influx database

First, you need to create a config for `influx` CLI:

```commandline
terraform output db_influx_config_cmd
```

Run the command that is output to set create the config.

Now you're ready to interact with InfluxDB, for example:

```console

```

## Testing

[webhook.site](https://webhook.site) is good for seeing what webhooks GitHub
sends us.

[ngrok](https://ngrok.com) is good for local testing of the backend.

[web application flow]: https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps#web-application-flow

[link_1]:
https://docs.github.com/en/apps/creating-github-apps/writing-code-for-a-github-app/building-ci-checks-with-a-github-app#authenticating-as-a-github-app

[login_btn]:
https://docs.github.com/en/apps/creating-github-apps/writing-code-for-a-github-app/building-a-login-with-github-button-with-a-github-app

[auth_gh_app]:
https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/authenticating-as-a-github-app

[auth_gh_install]:
https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/authenticating-as-a-github-app-installation

[auth_gh_user]:
https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/authenticating-with-a-github-app-on-behalf-of-a-user

[GitHub's token format]: https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/about-authentication-to-github#githubs-token-formats
