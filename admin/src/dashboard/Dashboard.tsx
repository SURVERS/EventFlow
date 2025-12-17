import React, { useEffect, useState } from 'react';
import { Card, CardContent, Typography, Grid, Box } from '@mui/material';
import {
  EventAvailable,
  People,
  ConfirmationNumber,
  Category,
  TrendingUp,
} from '@mui/icons-material';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
} from 'recharts';

interface DashboardStats {
  total_events: number;
  published_events: number;
  draft_events: number;
  total_participants: number;
  total_organizers: number;
  total_tickets: number;
  total_registrations: number;
  attended_count: number;
  attendance_rate: number;
}

interface PopularCategory {
  category_id: number;
  category_name: string;
  event_count: number;
}

const StatCard: React.FC<{
  title: string;
  value: number | string;
  icon: React.ReactNode;
  color: string;
}> = ({ title, value, icon, color }) => (
  <Card>
    <CardContent>
      <Box display="flex" alignItems="center" justifyContent="space-between">
        <Box>
          <Typography color="textSecondary" gutterBottom variant="body2">
            {title}
          </Typography>
          <Typography variant="h4">{value}</Typography>
        </Box>
        <Box
          sx={{
            backgroundColor: color,
            borderRadius: '50%',
            padding: 2,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
          }}
        >
          {icon}
        </Box>
      </Box>
    </CardContent>
  </Card>
);

export const Dashboard: React.FC = () => {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [popularCategories, setPopularCategories] = useState<PopularCategory[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch('http://localhost:8080/api/v1/dashboard/statistics')
      .then((response) => response.json())
      .then((data) => {
        setStats(data);
      })
      .catch((error) => console.error('Error fetching statistics:', error));

    fetch('http://localhost:8080/api/v1/dashboard/popular-categories')
      .then((response) => response.json())
      .then((data) => {
        setPopularCategories(data);
        setLoading(false);
      })
      .catch((error) => {
        console.error('Error fetching popular categories:', error);
        setLoading(false);
      });
  }, []);

  if (loading || !stats) {
    return <Typography>Загрузка...</Typography>;
  }

  return (
    <Box sx={{ padding: 3 }}>
      <Typography variant="h4" gutterBottom>
        Панель управления EventFlow
      </Typography>

      <Grid container spacing={3} sx={{ marginTop: 2 }}>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Всего событий"
            value={stats.total_events}
            icon={<EventAvailable sx={{ color: 'white' }} />}
            color="#3f51b5"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Участников"
            value={stats.total_participants}
            icon={<People sx={{ color: 'white' }} />}
            color="#4caf50"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Тикетов выдано"
            value={stats.total_tickets}
            icon={<ConfirmationNumber sx={{ color: 'white' }} />}
            color="#ff9800"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Посещаемость"
            value={`${stats.attendance_rate.toFixed(1)}%`}
            icon={<TrendingUp sx={{ color: 'white' }} />}
            color="#f44336"
          />
        </Grid>
      </Grid>

      <Grid container spacing={3} sx={{ marginTop: 2 }}>
        {/* График популярных категорий (Pie Chart) */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" marginBottom={2}>
                <Category sx={{ marginRight: 1 }} />
                <Typography variant="h6">Популярные категории</Typography>
              </Box>
              {popularCategories.length > 0 ? (
                <ResponsiveContainer width="100%" height={300}>
                  <PieChart>
                    <Pie
                      data={popularCategories}
                      dataKey="event_count"
                      nameKey="category_name"
                      cx="50%"
                      cy="50%"
                      outerRadius={80}
                      label={(entry) => `${entry.category_name}: ${entry.event_count}`}
                    >
                      {popularCategories.map((entry, index) => (
                        <Cell
                          key={`cell-${index}`}
                          fill={['#3f51b5', '#4caf50', '#ff9800', '#f44336', '#9c27b0'][index % 5]}
                        />
                      ))}
                    </Pie>
                    <Tooltip />
                  </PieChart>
                </ResponsiveContainer>
              ) : (
                <Typography variant="body2" color="textSecondary">
                  Нет данных о категориях
                </Typography>
              )}
            </CardContent>
          </Card>
        </Grid>

        {/* График статистики событий (Bar Chart) */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Статистика событий
              </Typography>
              <ResponsiveContainer width="100%" height={300}>
                <BarChart
                  data={[
                    { name: 'Опубликовано', value: stats.published_events },
                    { name: 'Черновики', value: stats.draft_events },
                    { name: 'Регистрации', value: stats.total_registrations },
                    { name: 'Посетило', value: stats.attended_count },
                  ]}
                >
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="name" />
                  <YAxis />
                  <Tooltip />
                  <Legend />
                  <Bar dataKey="value" fill="#667eea" />
                </BarChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
};
