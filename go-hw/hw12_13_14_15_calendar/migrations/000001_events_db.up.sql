CREATE TABLE IF NOT EXISTS events
(
    id          serial       NOT NULL,
    title       varchar(255) NOT NULL,
    description text         NULL,
    date_start  timestamptz  NOT NULL,
    date_end    timestamptz  NOT NULL,
    owner_id    int          NOT NULL,
    notify_in   int          NOT NULL,
    CONSTRAINT events_pk PRIMARY KEY (id)
);
CREATE UNIQUE INDEX owner_start_time_idx ON events USING btree (owner_id, date_start);
INSERT INTO events (id, title, description, date_start, date_end, owner_id, notify_in)
VALUES (999, 't', 'd', '2023-01-16T20:00:00Z', '2023-01-17T20:00:00Z', 1, 1),
       (1000, 't', 'd', '2023-01-16T20:00:00Z', '2023-01-17T20:00:00Z', 2, 1);