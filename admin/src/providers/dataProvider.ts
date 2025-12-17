import simpleRestProvider from 'ra-data-simple-rest';
import { fetchUtils } from 'react-admin';

const apiUrl = 'http://localhost:8080/api/v1';

const httpClient = (url: string, options: fetchUtils.Options = {}) => {
  const token = localStorage.getItem('access_token');

  if (!options.headers) {
    options.headers = new Headers({ Accept: 'application/json' });
  }

  if (token) {
    (options.headers as Headers).set('Authorization', `Bearer ${token}`);
  }

  return fetchUtils.fetchJson(url, options);
};

const dataProvider = simpleRestProvider(apiUrl, httpClient);

export default dataProvider;
