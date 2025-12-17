import React from 'react';
import { List, Datagrid, TextField, EmailField, EditButton, FunctionField, useRecordContext } from 'react-admin';
import { useEffect, useState } from 'react';

interface ParticipantStatistics {
  total_registrations: number;
  attendance_rate: number;
}

const StatisticsField: React.FC<{ source: string; label: string }> = ({ label }) => {
  const record = useRecordContext();
  const [statistics, setStatistics] = useState<ParticipantStatistics | null>(null);
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

  if (loading) return <span>...</span>;
  if (!statistics) return <span>-</span>;

  return <span>{statistics.total_registrations}</span>;
};

const AttendanceRateField: React.FC<{ source: string; label: string }> = ({ label }) => {
  const record = useRecordContext();
  const [statistics, setStatistics] = useState<ParticipantStatistics | null>(null);
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

  if (loading) return <span>...</span>;
  if (!statistics) return <span>-</span>;

  return <span>{statistics.attendance_rate.toFixed(1)}%</span>;
};

export const ParticipantList: React.FC = (props) => (
  <List {...props}>
    <Datagrid rowClick="edit">
      <TextField source="id" />
      <TextField source="full_name" label="ФИО" />
      <EmailField source="email" label="Email" />
      <TextField source="phone" label="Телефон" />
      <StatisticsField source="total_registrations" label="Регистраций" />
      <AttendanceRateField source="attendance_rate" label="Посещаемость" />
      <EditButton />
    </Datagrid>
  </List>
);
