import React from 'react';
import {
  Show,
  SimpleShowLayout,
  TextField,
  DateField,
  SelectField,
  ReferenceField,
  ReferenceManyField,
  Datagrid,
  TabbedShowLayout,
  Tab,
  NumberField,
} from 'react-admin';

export const EventShow: React.FC = (props) => (
  <Show {...props}>
    <TabbedShowLayout>
      <Tab label="Основная информация">
        <TextField source="id" />
        <TextField source="title" label="Название" />
        <TextField source="description" label="Описание" />
        <DateField source="start_time" label="Дата начала" showTime />
        <DateField source="end_time" label="Дата окончания" showTime />
        <TextField source="location" label="Место проведения" />
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
        <NumberField source="max_participants" label="Макс. участников" />
        <DateField source="created_at" label="Создано" showTime />
        <DateField source="updated_at" label="Обновлено" showTime />
      </Tab>

      <Tab label="Регистрации" path="registrations">
        <ReferenceManyField
          reference="event_registrations"
          target="event_id"
          label="Регистрации участников"
        >
          <Datagrid rowClick="show">
            <TextField source="id" />
            <ReferenceField source="participant_id" reference="participants" label="Участник">
              <TextField source="name" />
            </ReferenceField>
            <ReferenceField source="participant_id" reference="participants">
              <TextField source="email" label="Email" />
            </ReferenceField>
            <DateField source="registration_date" label="Дата регистрации" showTime />
            <SelectField
              source="status"
              label="Статус"
              choices={[
                { id: 'registered', name: 'Зарегистрирован' },
                { id: 'attended', name: 'Посетил' },
                { id: 'no-show', name: 'Не явился' },
              ]}
            />
          </Datagrid>
        </ReferenceManyField>
      </Tab>

      <Tab label="Тикеты" path="tickets">
        <ReferenceManyField reference="tickets" target="event_id" label="Выданные тикеты">
          <Datagrid rowClick="show">
            <TextField source="id" />
            <ReferenceField source="participant_id" reference="participants" label="Участник">
              <TextField source="name" />
            </ReferenceField>
            <TextField source="qr_code" label="QR-код" />
            <SelectField
              source="ticket_type"
              label="Тип"
              choices={[
                { id: 'free', name: 'Бесплатный' },
                { id: 'paid', name: 'Платный' },
              ]}
            />
            <SelectField
              source="status"
              label="Статус"
              choices={[
                { id: 'active', name: 'Активный' },
                { id: 'canceled', name: 'Отменен' },
              ]}
            />
            <DateField source="created_at" label="Выдан" showTime />
          </Datagrid>
        </ReferenceManyField>
      </Tab>
    </TabbedShowLayout>
  </Show>
);
