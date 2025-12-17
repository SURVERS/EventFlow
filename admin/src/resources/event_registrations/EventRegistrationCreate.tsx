import React from 'react';
import { Create, SimpleForm, ReferenceInput, SelectInput } from 'react-admin';

export const EventRegistrationCreate: React.FC = (props) => (
  <Create {...props}>
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
    </SimpleForm>
  </Create>
);
