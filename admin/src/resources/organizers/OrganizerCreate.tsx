import React from 'react';
import { Create, SimpleForm, TextInput, SelectInput } from 'react-admin';

export const OrganizerCreate: React.FC = (props) => (
  <Create {...props}>
    <SimpleForm>
      <TextInput source="name" label="Имя" fullWidth />
      <TextInput source="email" label="Email" type="email" fullWidth />
      <SelectInput
        source="role"
        label="Роль"
        choices={[
          { id: 'admin', name: 'Администратор' },
          { id: 'organizer', name: 'Организатор' },
        ]}
        fullWidth
      />
    </SimpleForm>
  </Create>
);
