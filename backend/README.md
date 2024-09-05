# backend

We try to keep them to a **reasonable** minimum.

## Auth

[GitHub Apps have three methods of authentication][link_1]. We support:

- [x] [Authenticating as a GitHub App][auth_gh_app] using a JSON Web Token (JWT)
- [x] [Authenticating as a specific installation][auth_gh_install] of a GitHub
  App using an installation access token
- [ ] [Authenticating on behalf of a user][auth_gh_user]

### How does auth work?

We use a thing called [web application flow].
It's also explained from the practical side in [Build a "Login" button tutorial][login_btn].

Overall it looks like this:

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

The access token lets us to make requests to the API on a behalf of a user.

For example, we can use it to access user's private repositories:

```console
curl \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  https://api.github.com/repos/bartekpacia/discover_rudy
```

This token can be entirely unrelated to the how sign in works in our webapp.

Our webapp can issue a JWT that is completely independent from the access token.

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

[web_appplication_flow]: https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps#web-application-flow
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
