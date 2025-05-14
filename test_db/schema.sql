DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS machines;
DROP TABLE IF EXISTS resources;
DROP TABLE IF EXISTS recipies;
DROP TABLE IF EXISTS recipies_inputs;
DROP TABLE IF EXISTS recipies_outputs;
DROP TABLE IF EXISTS machines_recipies;

CREATE TABLE users(
    id            integer PRIMARY KEY AUTOINCREMENT,
    name          text,
    password_hash text,
    role          text
);

CREATE TABLE machines(
    id             integer PRIMARY KEY AUTOINCREMENT,
    name           text,
    users_id       integer,
    inputs_solid   integer,
    inputs_liquid  integer,
    outputs_solid  integer,
    outputs_liquid integer,
    speed          integer,
    default_choice integer,
    FOREIGN KEY(users_id) REFERENCES users(id)
);

CREATE TABLE resources(
    id             integer PRIMARY KEY AUTOINCREMENT,
    name           text,
    users_id       integer,
    liquid         integer,
    resurce_unit   text,
    FOREIGN KEY(users_id) REFERENCES users(id)
);

CREATE TABLE recipies(
    id                    integer PRIMARY KEY AUTOINCREMENT,
    name                  text,
    users_id              integer,
    production_time       integer,
    production_time_unit  text,
    default_choice        integer,
    FOREIGN KEY(users_id) REFERENCES users(id)
);

CREATE TABLE recipies_inputs(
    id                    integer PRIMARY KEY AUTOINCREMENT,
    users_id              integer,
    recipies_id           integer,
    resources_id          integer,
    amount                integer,
    FOREIGN KEY(users_id) REFERENCES users(id),
    FOREIGN KEY(recipies_id) REFERENCES recipies(id),
    FOREIGN KEY(resources_id) REFERENCES resources(id)
);

CREATE TABLE recipies_outputs(
    id                    integer PRIMARY KEY AUTOINCREMENT,
    users_id              integer,
    recipies_id           integer,
    resources_id          integer,
    amount                integer,
    FOREIGN KEY(users_id) REFERENCES users(id),
    FOREIGN KEY(recipies_id) REFERENCES recipies(id),
    FOREIGN KEY(resources_id) REFERENCES resources(id)
);

CREATE TABLE machines_recipies(
    id                    integer PRIMARY KEY AUTOINCREMENT,
    users_id              integer,
    recipies_id           integer,
    machines_id           integer,
    FOREIGN KEY(users_id) REFERENCES users(id),
    FOREIGN KEY(recipies_id) REFERENCES recipies(id),
    FOREIGN KEY(machines_id) REFERENCES machines(id)
);