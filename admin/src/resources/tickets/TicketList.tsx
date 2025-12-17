import React from 'react';
import { List, Datagrid, TextField, ReferenceField, EditButton, FunctionField } from 'react-admin';

export const TicketList: React.FC = (props) => (
  <List {...props}>
    <Datagrid rowClick="edit">
      <TextField source="id" />
      <ReferenceField source="event_id" reference="events" label="Событие">
        <TextField source="title" />
      </ReferenceField>
      <ReferenceField source="participant_id" reference="participants" label="Участник">
        <TextField source="full_name" />
      </ReferenceField>
      <TextField source="ticket_type" label="Тип тикета" />
      <TextField source="status" label="Статус" />
      <FunctionField
        label="QR-код"
        render={(record: any) => (
          <span style={{ fontFamily: 'monospace', fontSize: '0.85em' }}>
            {record.qr_code ? record.qr_code.substring(0, 8) + '...' : '-'}
          </span>
        )}
      />
      <EditButton />
    </Datagrid>
  </List>
);
