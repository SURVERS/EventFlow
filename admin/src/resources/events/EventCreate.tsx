import React from 'react';
import { Create, SimpleForm, TextInput, DateTimeInput, SelectInput, ReferenceInput } from 'react-admin';

export const EventCreate: React.FC = (props) => (
  <Create {...props}>
    <SimpleForm>
      <TextInput source="title" label="Название" fullWidth />
      <TextInput source="description" label="Описание" multiline fullWidth />
      <DateTimeInput source="start_time" label="Дата и время начала" />
      <DateTimeInput source="end_time" label="Дата и время окончания" />
      <ReferenceInput source="category_id" reference="categories" label="Категория">
        <SelectInput optionText="name" fullWidth />
      </ReferenceInput>
      <SelectInput source="event_type" label="Тип события" choices={[
        { id: 'webinar', name: 'Вебинар' },
        { id: 'workshop', name: 'Воркшоп' },
        { id: 'conference', name: 'Конференция' }
      ]} fullWidth />
      <SelectInput source="status" label="Статус события" choices={[
        { id: 'scheduled', name: 'Запланировано' },
        { id: 'ongoing', name: 'В процессе' },
        { id: 'completed', name: 'Завершено' },
        { id: 'canceled', name: 'Отменено' }
      ]} defaultValue="scheduled" fullWidth />
      <SelectInput source="publish_status" label="Статус публикации" choices={[
        { id: 'draft', name: 'Черновик' },
        { id: 'published', name: 'Опубликовано' }
      ]} defaultValue="draft" fullWidth />
    </SimpleForm>
  </Create>
);