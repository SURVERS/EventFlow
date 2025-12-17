import type { AuthProvider } from 'react-admin';

const API_URL = 'http://localhost:8080/api/v1';

let refreshTokenTimeout: NodeJS.Timeout | null = null;

const scheduleTokenRefresh = (expiresIn: number) => {
  if (refreshTokenTimeout) {
    clearTimeout(refreshTokenTimeout);
  }

  const refreshTime = (expiresIn - 30) * 1000;

  refreshTokenTimeout = setTimeout(async () => {
    try {
      const refreshToken = localStorage.getItem('refresh_token');
      if (!refreshToken) {
        return;
      }

      const response = await fetch(`${API_URL}/auth/refresh`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refresh_token: refreshToken }),
      });

      if (response.ok) {
        const data = await response.json();
        localStorage.setItem('access_token', data.access_token);
        localStorage.setItem('user', JSON.stringify(data.user));
        scheduleTokenRefresh(data.expires_in);
      } else {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('user');
        window.location.href = '/login';
      }
    } catch (error) {
      console.error('Token refresh failed:', error);
    }
  }, refreshTime);
};

export const authProvider: AuthProvider = {
  login: async ({ username, password }) => {
    const request = new Request(`${API_URL}/auth/login`, {
      method: 'POST',
      body: JSON.stringify({ email: username, password }),
      headers: new Headers({ 'Content-Type': 'application/json' }),
    });

    const response = await fetch(request);
    if (response.status < 200 || response.status >= 300) {
      throw new Error('Invalid email or password');
    }

    const data = await response.json();

    localStorage.setItem('access_token', data.access_token);
    localStorage.setItem('refresh_token', data.refresh_token);
    localStorage.setItem('user', JSON.stringify(data.user));

    scheduleTokenRefresh(data.expires_in);

    return Promise.resolve();
  },

  logout: () => {
    if (refreshTokenTimeout) {
      clearTimeout(refreshTokenTimeout);
    }
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user');
    return Promise.resolve();
  },

  checkAuth: () => {
    const token = localStorage.getItem('access_token');
    if (!token) {
      return Promise.reject();
    }

    const refreshToken = localStorage.getItem('refresh_token');
    if (refreshToken && !refreshTokenTimeout) {
      scheduleTokenRefresh(300);
    }

    return Promise.resolve();
  },

  checkError: (error) => {
    const status = error.status;
    if (status === 401 || status === 403) {
      if (refreshTokenTimeout) {
        clearTimeout(refreshTokenTimeout);
      }
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      localStorage.removeItem('user');
      return Promise.reject();
    }
    return Promise.resolve();
  },

  getPermissions: () => {
    const userStr = localStorage.getItem('user');
    if (!userStr) {
      return Promise.reject();
    }

    const user = JSON.parse(userStr);
    return Promise.resolve(user.role);
  },

  getIdentity: () => {
    const userStr = localStorage.getItem('user');
    if (!userStr) {
      return Promise.reject();
    }

    const user = JSON.parse(userStr);
    return Promise.resolve({
      id: user.id,
      fullName: user.name,
      avatar: undefined,
    });
  },
};
