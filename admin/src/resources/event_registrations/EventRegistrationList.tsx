import React from 'react';
import { List, Datagrid, TextField, DateField, ReferenceField, EditButton } from 'react-admin';

export const EventRegistrationList: React.FC = (props) => (
  <List {...props}>
    <Datagrid rowClick="edit">
      <TextField source="id" />
      <ReferenceField source="event_id" reference="events" label="Событие">
        <TextField source="title" />
      </ReferenceField>
      <ReferenceField source="participant_id" reference="participants" label="Участник">
        <TextField source="full_name" />
      </ReferenceField>
      <TextField source="status" label="Статус" />
      <DateField source="registered_at" label="Дата регистрации" showTime />
      <EditButton />
    </Datagrid>
  </List>
);
