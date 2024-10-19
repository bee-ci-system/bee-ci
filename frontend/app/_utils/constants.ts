function loadApiBaseUrl(): string {
  const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL;
  if (!apiBaseUrl) {
    throw new Error('NEXT_PUBLIC_API_BASE_URL is not set');
  }
  return apiBaseUrl;
}

function loadGithubAppClientId(): string {
  const githubAppClientId = process.env.NEXT_PUBLIC_GITHUB_APP_CLIENT_ID;
  if (!githubAppClientId) {
    throw new Error('NEXT_PUBLIC_GITHUB_APP_CLIENT_ID is not set');
  }
  return githubAppClientId;
}

export const apiBaseUrl = loadApiBaseUrl();
export const authUrl = `https://github.com/login/oauth/authorize?client_id=${loadGithubAppClientId()}`;
