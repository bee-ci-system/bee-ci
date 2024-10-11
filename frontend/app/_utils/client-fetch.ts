import Cookies from 'js-cookie';

export const clientFetch = async (
  endpoint: string,
  options: RequestInit = {},
) => {
  const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL;

  if (!baseUrl) {
    throw new Error('Base url is not established');
  }

  const token = Cookies.get('jwt');

  const headers: HeadersInit = new Headers({
    'Content-Type': 'application/json',
    'Access-Control-Allow-Origin': '*',
    ...options.headers,
  });

  if (token) {
    headers.set('Authorization', `bearer ${token}`);
  }

  const response = await fetch(`${baseUrl}${endpoint}`, {
    ...options,
    headers,
    cache: 'no-store',
  });

  return response;
};
