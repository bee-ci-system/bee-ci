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

  const response = await fetch(`${apiBaseUrl}${endpoint}`, {
    ...options,
    headers,
    cache: 'no-store',
  });

  return response;
};
