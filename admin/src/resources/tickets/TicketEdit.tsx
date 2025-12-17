import React from 'react';
import { Edit, SelectInput, TextField, ReferenceField, TabbedForm, FormTab, useRecordContext } from 'react-admin';
import { Box, Typography } from '@mui/material';
import { QRCodeSVG } from 'qrcode.react';

const QRCodeDisplay: React.FC = () => {
  const record = useRecordContext();

  if (!record || !record.qr_code) {
    return <Typography>QR-код отсутствует</Typography>;
  }

  return (
    <Box sx={{ padding: 2, textAlign: 'center' }}>
      <Typography variant="h6" gutterBottom>
        QR-код тикета
      </Typography>
      <Box sx={{ background: 'white', padding: 2, display: 'inline-block', marginBottom: 2 }}>
        <QRCodeSVG value={record.qr_code} size={200} />
      </Box>
      <Typography variant="body2" sx={{ fontFamily: 'monospace', wordBreak: 'break-all' }}>
        {record.qr_code}
      </Typography>
    </Box>
  );
};

export const TicketEdit: React.FC = (props) => (
  <Edit {...props}>
    <TabbedForm>
      <FormTab label="Основная информация">
        <ReferenceField source="event_id" reference="events" label="Событие">
          <TextField source="title" />
        </ReferenceField>
        <ReferenceField source="participant_id" reference="participants" label="Участник">
          <TextField source="full_name" />
        </ReferenceField>
        <SelectInput
          source="ticket_type"
          label="Тип тикета"
          choices={[
            { id: 'free', name: 'Бесплатный' },
            { id: 'paid', name: 'Платный' },
          ]}
          fullWidth
        />
        <SelectInput
          source="status"
          label="Статус"
          choices={[
            { id: 'active', name: 'Активный' },
            { id: 'canceled', name: 'Отменён' },
          ]}
          fullWidth
        />
      </FormTab>
      <FormTab label="QR-код">
        <TextField source="qr_code" label="QR-код" fullWidth disabled />
        <QRCodeDisplay />
      </FormTab>
    </TabbedForm>
  </Edit>
);
