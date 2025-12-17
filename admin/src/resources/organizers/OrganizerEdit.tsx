import React from 'react';
import { Edit, SimpleForm, TextInput, SelectInput } from 'react-admin';

export const OrganizerEdit: React.FC = (props) => (
  <Edit {...props}>
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
  </Edit>
);
