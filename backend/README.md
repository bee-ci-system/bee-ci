# backend

### Auth

[GitHub Apps have three methods of authentication][link_1]. We support:

- [x] [Authenticating as a GitHub App][auth_gh_app] using a JSON Web Token (JWT)
- [x] [Authenticating as a specific installation][auth_gh_install] of a GitHub App using an
  installation access token
- [ ] [Authenticating on behalf of a user][auth_gh_user]


### Testing

[webhook.site](https://webhook.site) is good for seeing what webhooks GitHub sends us.

[ngrok](https://ngrok.com) is good for local testing of the backend.

[link_1]: https://docs.github.com/en/apps/creating-github-apps/writing-code-for-a-github-app/building-ci-checks-with-a-github-app#authenticating-as-a-github-app


[auth_gh_app]: https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/authenticating-as-a-github-app
[auth_gh_install]: https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/authenticating-as-a-github-app-installation
[auth_gh_user]: https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/authenticating-with-a-github-app-on-behalf-of-a-user