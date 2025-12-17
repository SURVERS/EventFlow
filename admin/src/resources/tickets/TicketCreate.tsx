import React from 'react';
import { Create, SimpleForm, ReferenceInput, SelectInput } from 'react-admin';

export const TicketCreate: React.FC = (props) => (
  <Create {...props}>
    <SimpleForm>
      <ReferenceInput source="event_id" reference="events" label="Событие">
        <SelectInput optionText="title" fullWidth />
      </ReferenceInput>
      <ReferenceInput source="participant_id" reference="participants" label="Участник">
        <SelectInput optionText="full_name" fullWidth />
      </ReferenceInput>
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
    </SimpleForm>
  </Create>
);
