import React, { useEffect, useState } from 'react';
import { Edit, SimpleForm, TextInput, TabbedForm, FormTab, useRecordContext } from 'react-admin';
import { Card, CardContent, Typography, Grid, Box } from '@mui/material';

interface DetailedStatistics {
  participant_id: number;
  full_name: string;
  total_registrations: number;
  registered_count: number;
  attended_count: number;
  no_show_count: number;
  attendance_rate: number;
}

const StatisticsPanel: React.FC = () => {
  const record = useRecordContext();
  const [statistics, setStatistics] = useState<DetailedStatistics | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (record && record.id) {
      fetch(`http://localhost:8080/api/v1/participants/${record.id}/statistics`)
        .then(response => response.json())
        .then(data => {
          setStatistics(data);
          setLoading(false);
        })
        .catch(error => {
          console.error('Error fetching statistics:', error);
          setLoading(false);
        });
    }
  }, [record]);

  if (loading) {
    return <Typography>Загрузка статистики...</Typography>;
  }

  if (!statistics) {
    return <Typography>Не удалось загрузить статистику</Typography>;
  }

  return (
    <Box sx={{ padding: 2 }}>
      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Общая информация
              </Typography>
              <Typography variant="body1">
                <strong>Участник:</strong> {statistics.full_name}
              </Typography>
              <Typography variant="body1">
                <strong>Всего регистраций:</strong> {statistics.total_registrations}
              </Typography>
              <Typography variant="body1">
                <strong>Процент посещаемости:</strong> {statistics.attendance_rate.toFixed(1)}%
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Разбивка по статусам
              </Typography>
              <Typography variant="body1">
                <strong>Зарегистрирован:</strong> {statistics.registered_count}
              </Typography>
              <Typography variant="body1">
                <strong>Посетил:</strong> {statistics.attended_count}
              </Typography>
              <Typography variant="body1">
                <strong>Не пришёл:</strong> {statistics.no_show_count}
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Анализ активности
              </Typography>
              {statistics.total_registrations === 0 ? (
                <Typography variant="body1">
                  Участник ещё не регистрировался на мероприятия
                </Typography>
              ) : (
                <>
                  <Typography variant="body1">
                    Участник зарегистрировался на <strong>{statistics.total_registrations}</strong> мероприятий(я).
                  </Typography>
                  <Typography variant="body1">
                    Из них посетил <strong>{statistics.attended_count}</strong> ({statistics.attendance_rate.toFixed(1)}%).
                  </Typography>
                  {statistics.no_show_count > 0 && (
                    <Typography variant="body1" color="error">
                      Не пришёл на <strong>{statistics.no_show_count}</strong> мероприятий(я).
                    </Typography>
                  )}
                </>
              )}
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
};

export const ParticipantEdit: React.FC = (props) => (
  <Edit {...props}>
    <TabbedForm>
      <FormTab label="Основная информация">
        <TextInput source="full_name" label="ФИО" fullWidth />
        <TextInput source="email" label="Email" type="email" fullWidth />
        <TextInput source="phone" label="Телефон" fullWidth />
      </FormTab>
      <FormTab label="Статистика">
        <StatisticsPanel />
      </FormTab>
    </TabbedForm>
  </Edit>
);
