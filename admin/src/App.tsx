import { Admin, Resource } from 'react-admin';
import dataProvider from './providers/dataProvider';
import { authProvider } from './providers/authProvider';
import { Login } from './components/Login';
import { Dashboard } from './dashboard/Dashboard';
import { CategoryList } from './resources/categories/CategoryList';
import { CategoryCreate } from './resources/categories/CategoryCreate';
import { CategoryEdit } from './resources/categories/CategoryEdit';
import { EventList } from './resources/events/EventList';
import { EventCreate } from './resources/events/EventCreate';
import { EventEdit } from './resources/events/EventEdit';
import { EventShow } from './resources/events/EventShow';
import { ParticipantList } from './resources/participants/ParticipantList';
import { ParticipantCreate } from './resources/participants/ParticipantCreate';
import { ParticipantEdit } from './resources/participants/ParticipantEdit';
import { EventRegistrationCreate } from './resources/event_registrations/EventRegistrationCreate';
import { EventRegistrationEdit } from './resources/event_registrations/EventRegistrationEdit';
import { EventRegistrationList } from './resources/event_registrations/EventRegistrationList';
import { OrganizerList } from './resources/organizers/OrganizerList';
import { OrganizerCreate } from './resources/organizers/OrganizerCreate';
import { OrganizerEdit } from './resources/organizers/OrganizerEdit';
import { TicketList } from './resources/tickets/TicketList';
import { TicketCreate } from './resources/tickets/TicketCreate';
import { TicketEdit } from './resources/tickets/TicketEdit';

const App = () => (
    <Admin
        dataProvider={dataProvider}
        authProvider={authProvider}
        loginPage={Login}
        dashboard={Dashboard}
    >
        <Resource 
            name="categories" 
            list={CategoryList}
            edit={CategoryEdit}
            create={CategoryCreate}
        />
        <Resource
            name="events"
            list={EventList}
            edit={EventEdit}
            create={EventCreate}
            show={EventShow}
        />
        <Resource 
            name="participants" 
            list={ParticipantList}
            edit={ParticipantEdit}
            create={ParticipantCreate}
        />
        <Resource
            name="event_registrations"
            list={EventRegistrationList}
            edit={EventRegistrationEdit}
            create={EventRegistrationCreate}
        />
        <Resource
            name="organizers"
            list={OrganizerList}
            edit={OrganizerEdit}
            create={OrganizerCreate}
        />
        <Resource
            name="tickets"
            list={TicketList}
            edit={TicketEdit}
            create={TicketCreate}
        />
    </Admin>
);

export default App;