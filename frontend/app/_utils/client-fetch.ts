import Cookies from 'js-cookie';
import { apiBaseUrl } from '../_utils/constants';

export const clientFetch = async (
  endpoint: string,
  options: RequestInit = {},
) => {
  const token = Cookies.get('jwt');

  const headers: HeadersInit = new Headers({
    'Content-Type': 'application/json',
    'Access-Control-Allow-Origin': '*',
    ...options.headers,
  });

  if (token) {
    headers.set('Authorization', `bearer ${token}`);
  }

  const url = `${apiBaseUrl}${endpoint}`;
  console.log(`clientFetch: will fetch from: ${url}`);
  const response = await fetch(url, {
    ...options,
    headers,
    cache: 'no-store',
  });

  return response;
};
