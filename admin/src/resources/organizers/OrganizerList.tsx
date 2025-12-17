import React from 'react';
import { List, Datagrid, TextField, EmailField, EditButton } from 'react-admin';

export const OrganizerList: React.FC = (props) => (
  <List {...props}>
    <Datagrid rowClick="edit">
      <TextField source="id" />
      <TextField source="name" label="Имя" />
      <EmailField source="email" label="Email" />
      <TextField source="role" label="Роль" />
      <EditButton />
    </Datagrid>
  </List>
);
