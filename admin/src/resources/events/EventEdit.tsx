import React from 'react';
import { Edit, SimpleForm, TextInput, DateTimeInput, SelectInput, ReferenceInput } from 'react-admin';

const eventTypeChoices = [
  { id: 'webinar', name: 'Вебинар' },
  { id: 'workshop', name: 'Воркшоп' },
  { id: 'conference', name: 'Конференция' }
];

const statusChoices = [
  { id: 'scheduled', name: 'Запланировано' },
  { id: 'ongoing', name: 'В процессе' },
  { id: 'completed', name: 'Завершено' },
  { id: 'canceled', name: 'Отменено' }
];

const publishStatusChoices = [
  { id: 'draft', name: 'Черновик' },
  { id: 'published', name: 'Опубликовано' }
];

export const EventEdit: React.FC = (props) => (
  <Edit {...props}>
    <SimpleForm>
      <TextInput source="title" label="Название" fullWidth />
      <TextInput source="description" label="Описание" multiline fullWidth />
      <DateTimeInput source="start_time" label="Дата и время начала" />
      <DateTimeInput source="end_time" label="Дата и время окончания" />
      <ReferenceInput source="category_id" reference="categories" label="Категория">
        <SelectInput optionText="name" fullWidth />
      </ReferenceInput>
      <SelectInput source="event_type" label="Тип события" choices={eventTypeChoices} fullWidth />
      <SelectInput source="status" label="Статус события" choices={statusChoices} fullWidth />
      <SelectInput source="publish_status" label="Статус публикации" choices={publishStatusChoices} fullWidth />
    </SimpleForm>
  </Edit>
);
