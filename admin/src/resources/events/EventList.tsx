import React from 'react';
import {
  List,
  Datagrid,
  TextField,
  DateField,
  SelectField,
  EditButton,
  ReferenceField,
  TextInput,
  SelectInput,
  ReferenceInput,
  DateInput,
  ShowButton,
} from 'react-admin';

const eventFilters = [
  <TextInput key="q" label="Поиск" source="q" alwaysOn />,
  <ReferenceInput key="category_id" source="category_id" reference="categories" label="Категория">
    <SelectInput optionText="name" />
  </ReferenceInput>,
  <SelectInput
    key="event_type"
    source="event_type"
    label="Тип"
    choices={[
      { id: 'webinar', name: 'Вебинар' },
      { id: 'workshop', name: 'Воркшоп' },
      { id: 'conference', name: 'Конференция' },
    ]}
  />,
  <SelectInput
    key="status"
    source="status"
    label="Статус"
    choices={[
      { id: 'scheduled', name: 'Запланировано' },
      { id: 'ongoing', name: 'В процессе' },
      { id: 'completed', name: 'Завершено' },
      { id: 'canceled', name: 'Отменено' },
    ]}
  />,
  <SelectInput
    key="publish_status"
    source="publish_status"
    label="Публикация"
    choices={[
      { id: 'draft', name: 'Черновик' },
      { id: 'published', name: 'Опубликовано' },
    ]}
  />,
  <DateInput key="start_date" source="start_date" label="Дата начала от" />,
  <DateInput key="end_date" source="end_date" label="Дата окончания до" />,
];

export const EventList: React.FC = (props) => (
  <List {...props} filters={eventFilters}>
    <Datagrid rowClick="edit">
      <TextField source="id" />
      <TextField source="title" label="Название" />
      <DateField source="start_time" label="Дата начала" showTime />
      <ReferenceField source="category_id" reference="categories" label="Категория">
        <TextField source="name" />
      </ReferenceField>
      <SelectField
        source="event_type"
        label="Тип"
        choices={[
          { id: 'webinar', name: 'Вебинар' },
          { id: 'workshop', name: 'Воркшоп' },
          { id: 'conference', name: 'Конференция' },
        ]}
      />
      <SelectField
        source="status"
        label="Статус"
        choices={[
          { id: 'scheduled', name: 'Запланировано' },
          { id: 'ongoing', name: 'В процессе' },
          { id: 'completed', name: 'Завершено' },
          { id: 'canceled', name: 'Отменено' },
        ]}
      />
      <SelectField
        source="publish_status"
        label="Публикация"
        choices={[
          { id: 'draft', name: 'Черновик' },
          { id: 'published', name: 'Опубликовано' },
        ]}
      />
      <ShowButton />
      <EditButton />
    </Datagrid>
  </List>
);