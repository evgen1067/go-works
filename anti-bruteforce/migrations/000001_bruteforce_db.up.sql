CREATE TABLE IF NOT EXISTS blacklist
(
    id      serial      NOT NULL PRIMARY KEY,
    address inet        NOT NULL UNIQUE,
    created timestamptz NOT NULL default now()
);
CREATE INDEX ON blacklist USING GIST (address inet_ops);

CREATE TABLE IF NOT EXISTS whitelist
(
    id      serial      NOT NULL PRIMARY KEY,
    address inet        NOT NULL UNIQUE,
    created timestamptz NOT NULL default now()
);
CREATE INDEX ON whitelist USING GIST (address inet_ops);