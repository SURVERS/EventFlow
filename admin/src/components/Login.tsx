import React, { useState } from 'react';
import { useLogin, useNotify } from 'react-admin';
import {
  Box,
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  Tabs,
  Tab,
  Alert,
  CircularProgress,
} from '@mui/material';
import { EventAvailable, Login as LoginIcon, PersonAdd } from '@mui/icons-material';

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

const TabPanel = (props: TabPanelProps) => {
  const { children, value, index, ...other } = props;
  return (
    <div role="tabpanel" hidden={value !== index} {...other}>
      {value === index && <Box sx={{ pt: 3 }}>{children}</Box>}
    </div>
  );
};

export const Login = () => {
  const [tabValue, setTabValue] = useState(0);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [name, setName] = useState('');
  const [role, setRole] = useState('organizer');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const login = useLogin();
  const notify = useNotify();

  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
    setError('');
  };

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      await login({ username: email, password });
    } catch (err) {
      setError('Неверный email или пароль');
      setLoading(false);
    }
  };

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const response = await fetch('http://localhost:8080/api/v1/auth/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, email, password, role }),
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.error || 'Ошибка регистрации');
      }

      const data = await response.json();
      localStorage.setItem('token', data.token);
      localStorage.setItem('user', JSON.stringify(data.user));

      notify('Регистрация успешна!', { type: 'success' });
      window.location.href = '/';
    } catch (err: any) {
      setError(err.message || 'Ошибка при регистрации');
      setLoading(false);
    }
  };

  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'column',
        minHeight: '100vh',
        alignItems: 'center',
        justifyContent: 'center',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        padding: 2,
      }}
    >
      <Card
        sx={{
          minWidth: { xs: '100%', sm: 450 },
          maxWidth: 500,
          boxShadow: '0 8px 32px rgba(0, 0, 0, 0.3)',
          borderRadius: 3,
        }}
      >
        <CardContent sx={{ padding: 4 }}>
          {/* Заголовок */}
          <Box sx={{ textAlign: 'center', marginBottom: 3 }}>
            <EventAvailable
              sx={{
                fontSize: 60,
                color: '#667eea',
                marginBottom: 1,
              }}
            />
            <Typography variant="h4" gutterBottom fontWeight="bold">
              EventFlow
            </Typography>
            <Typography variant="body2" color="textSecondary">
              Система управления событиями
            </Typography>
          </Box>

          {/* Табы Вход/Регистрация */}
          <Tabs
            value={tabValue}
            onChange={handleTabChange}
            variant="fullWidth"
            sx={{ borderBottom: 1, borderColor: 'divider', marginBottom: 2 }}
          >
            <Tab icon={<LoginIcon />} label="Вход" />
            <Tab icon={<PersonAdd />} label="Регистрация" />
          </Tabs>

          {/* Панель входа */}
          <TabPanel value={tabValue} index={0}>
            <form onSubmit={handleLogin}>
              {error && (
                <Alert severity="error" sx={{ marginBottom: 2 }}>
                  {error}
                </Alert>
              )}
              <TextField
                label="Email"
                type="email"
                fullWidth
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                margin="normal"
                variant="outlined"
              />
              <TextField
                label="Пароль"
                type="password"
                fullWidth
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                margin="normal"
                variant="outlined"
              />
              <Button
                type="submit"
                fullWidth
                variant="contained"
                size="large"
                disabled={loading}
                sx={{
                  marginTop: 3,
                  padding: 1.5,
                  background: 'linear-gradient(45deg, #667eea 30%, #764ba2 90%)',
                  '&:hover': {
                    background: 'linear-gradient(45deg, #5568d3 30%, #63408a 90%)',
                  },
                }}
              >
                {loading ? <CircularProgress size={24} color="inherit" /> : 'Войти'}
              </Button>
            </form>
          </TabPanel>

          {/* Панель регистрации */}
          <TabPanel value={tabValue} index={1}>
            <form onSubmit={handleRegister}>
              {error && (
                <Alert severity="error" sx={{ marginBottom: 2 }}>
                  {error}
                </Alert>
              )}
              <TextField
                label="Имя"
                fullWidth
                required
                value={name}
                onChange={(e) => setName(e.target.value)}
                margin="normal"
                variant="outlined"
              />
              <TextField
                label="Email"
                type="email"
                fullWidth
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                margin="normal"
                variant="outlined"
              />
              <TextField
                label="Пароль"
                type="password"
                fullWidth
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                margin="normal"
                variant="outlined"
                helperText="Минимум 6 символов"
              />
              <TextField
                select
                label="Роль"
                fullWidth
                required
                value={role}
                onChange={(e) => setRole(e.target.value)}
                margin="normal"
                variant="outlined"
                SelectProps={{
                  native: true,
                }}
              >
                <option value="organizer">Организатор</option>
                <option value="admin">Администратор</option>
              </TextField>
              <Button
                type="submit"
                fullWidth
                variant="contained"
                size="large"
                disabled={loading}
                sx={{
                  marginTop: 3,
                  padding: 1.5,
                  background: 'linear-gradient(45deg, #667eea 30%, #764ba2 90%)',
                  '&:hover': {
                    background: 'linear-gradient(45deg, #5568d3 30%, #63408a 90%)',
                  },
                }}
              >
                {loading ? <CircularProgress size={24} color="inherit" /> : 'Зарегистрироваться'}
              </Button>
            </form>
          </TabPanel>

          {/* Подсказка */}
          <Box sx={{ marginTop: 3, textAlign: 'center' }}>
            <Typography variant="caption" color="textSecondary">
              {tabValue === 0
                ? 'Нет аккаунта? Перейдите на вкладку "Регистрация"'
                : 'Уже есть аккаунт? Перейдите на вкладку "Вход"'}
            </Typography>
          </Box>
        </CardContent>
      </Card>

      {/* Футер */}
      <Typography
        variant="body2"
        sx={{ marginTop: 4, color: 'white', opacity: 0.8 }}
      >
        © 2025 EventFlow. Все права защищены.
      </Typography>
    </Box>
  );
};
