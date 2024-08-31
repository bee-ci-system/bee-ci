# backend

We try to keep them to a **reasonable** minimum.

## Auth

[GitHub Apps have three methods of authentication][link_1]. We support:

- [x] [Authenticating as a GitHub App][auth_gh_app] using a JSON Web Token (JWT)
- [x] [Authenticating as a specific installation][auth_gh_install] of a GitHub
  App using an installation access token
- [ ] [Authenticating on behalf of a user][auth_gh_user]

### How does auth work?

The process is explained in the [Build a "Login" button tutorial][login_btn].

Overall it looks like this:

- The user clicks the "Sign in with GitHub" button and is redirected
  `https://github.com/login/oauth/authorize?client_id=$CLIENT_ID`

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

  GitHub responds with an access token.



## Database

We don't use an ORM, instead we use [sqlx](https://jmoiron.github.io/sqlx).

[pgcli](https://www.pgcli.com) is a good client for PostgreSQL.

```console
pgcli -h localhost -p 5432 -u postgres -W -d bee
```

List schemas:

```postgresql
\dn
```

List tables:

```postgresql
\dt bee_schema.*
```

Example query:

```postgresql
SELECT * FROM bee_schema.users
```

### Testing

[webhook.site](https://webhook.site) is good for seeing what webhooks GitHub
sends us.

[ngrok](https://ngrok.com) is good for local testing of the backend.

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
