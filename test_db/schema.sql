DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS machines;
DROP TABLE IF EXISTS resources;
DROP TABLE IF EXISTS recipes;
DROP TABLE IF EXISTS recipes_inputs;
DROP TABLE IF EXISTS recipes_outputs;
DROP TABLE IF EXISTS machines_recipes;

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
    speed          real,
    power_consumption_kw integer,
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

CREATE TABLE recipes(
    id                    integer PRIMARY KEY AUTOINCREMENT,
    name                  text,
    users_id              integer,
    production_time_s     integer,
    default_choice        integer,
    FOREIGN KEY(users_id) REFERENCES users(id)
);

CREATE TABLE recipes_inputs(
    id                    integer PRIMARY KEY AUTOINCREMENT,
    users_id              integer,
    recipes_id            integer,
    resources_id          integer,
    amount                integer,
    FOREIGN KEY(users_id) REFERENCES users(id),
    FOREIGN KEY(recipes_id) REFERENCES recipes(id),
    FOREIGN KEY(resources_id) REFERENCES resources(id)
);

CREATE TABLE recipes_outputs(
    id                    integer PRIMARY KEY AUTOINCREMENT,
    users_id              integer,
    recipes_id            integer,
    resources_id          integer,
    amount                integer,
    FOREIGN KEY(users_id) REFERENCES users(id),
    FOREIGN KEY(recipes_id) REFERENCES recipes(id),
    FOREIGN KEY(resources_id) REFERENCES resources(id)
);

CREATE TABLE machines_recipes(
    id                    integer PRIMARY KEY AUTOINCREMENT,
    users_id              integer,
    recipes_id           integer,
    machines_id           integer,
    FOREIGN KEY(users_id) REFERENCES users(id),
    FOREIGN KEY(recipes_id) REFERENCES recipes(id),
    FOREIGN KEY(machines_id) REFERENCES machines(id)
);