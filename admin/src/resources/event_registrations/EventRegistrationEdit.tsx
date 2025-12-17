import React from 'react';
import { Edit, SimpleForm, ReferenceInput, SelectInput, DateField } from 'react-admin';

export const EventRegistrationEdit: React.FC = (props) => (
  <Edit {...props}>
    <SimpleForm>
      <ReferenceInput source="event_id" reference="events" label="Событие">
        <SelectInput optionText="title" />
      </ReferenceInput>
      <ReferenceInput source="participant_id" reference="participants" label="Участник">
        <SelectInput optionText="full_name" />
      </ReferenceInput>
      <SelectInput source="status" label="Статус" choices={[
        { id: 'registered', name: 'Зарегистрирован' },
        { id: 'attended', name: 'Присутствовал' },
        { id: 'no-show', name: 'Не пришел' }
      ]} />
      <DateField source="registered_at" label="Дата регистрации" showTime disabled />
    </SimpleForm>
  </Edit>
);
