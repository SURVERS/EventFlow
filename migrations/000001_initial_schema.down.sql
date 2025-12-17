-- Rollback initial schema migration
-- Drops all tables in reverse order (respecting foreign keys)

DROP TABLE IF EXISTS tickets;
DROP TABLE IF EXISTS event_registrations;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS participants;
DROP TABLE IF EXISTS organizers;
DROP TABLE IF EXISTS event_types;
DROP TABLE IF EXISTS categories;
