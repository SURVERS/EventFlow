import React from 'react';
import { Create, SimpleForm, TextInput } from 'react-admin';

export const ParticipantCreate: React.FC = (props) => (
  <Create {...props}>
    <SimpleForm>
      <TextInput source="full_name" label="ФИО" fullWidth />
      <TextInput source="email" label="Email" type="email" fullWidth />
      <TextInput source="phone" label="Телефон" fullWidth />
    </SimpleForm>
  </Create>
);
