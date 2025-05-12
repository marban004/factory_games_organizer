DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS machines;
DROP TABLE IF EXISTS resources;
DROP TABLE IF EXISTS recipies;
DROP TABLE IF EXISTS recipies_inputs;
DROP TABLE IF EXISTS recipies_outputs;
DROP TABLE IF EXISTS machines_recipies;

CREATE TABLE users(
    id            integer,
    name          text,
    password_hash text,
    role          text
);

CREATE TABLE machines(
    id             integer,
    name           text,
    users_id       integer,
    inputs_solid   integer,
    inputs_liquid  integer,
    outputs_solid  integer,
    outputs_liquid integer,
    speed          integer,
    default_choice integer
);

CREATE TABLE resources(
    id             integer,
    name           text,
    users_id       integer,
    liquid         integer,
    resurce_unit   text
);

CREATE TABLE recipies(
    id                    integer,
    name                  text,
    users_id              integer,
    production_time       integer,
    production_time_unit  text,
    default_choice        integer
);

CREATE TABLE recipies_inputs(
    id                    integer,
    users_id              integer,
    recipies_id           integer,
    resources_id          integer,
    amount                integer
);

CREATE TABLE recipies_outputs(
    id                    integer,
    users_id              integer,
    recipies_id           integer,
    resources_id          integer,
    amount                integer
);

CREATE TABLE machines_recipies(
    id                    integer,
    users_id              integer,
    recipies_id           integer,
    machines_id           integer
);