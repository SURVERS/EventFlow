import { List, Datagrid, TextField, DateField } from 'react-admin';

export const CategoryList = () => (
    <List>
        <Datagrid rowClick="edit"> 
            <TextField source="id" />
            <TextField source="name" />
            <DateField source="created_at" showTime /> 
            <DateField source="updated_at" showTime />
        </Datagrid>
    </List>
);