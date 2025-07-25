DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS machines_recipes;
DROP TABLE IF EXISTS recipes_outputs;
DROP TABLE IF EXISTS recipes_inputs;
DROP TABLE IF EXISTS recipes;
DROP TABLE IF EXISTS resources;
DROP TABLE IF EXISTS machines;

CREATE TABLE machines(
    id             integer PRIMARY KEY AUTO_INCREMENT,
    name           text,
    users_id       integer,
    inputs_solid   integer,
    inputs_liquid  integer,
    outputs_solid  integer,
    outputs_liquid integer,
    speed          real,
    power_consumption_kw integer,
    default_choice integer
);

CREATE TABLE resources(
    id              integer PRIMARY KEY AUTO_INCREMENT,
    name            text,
    users_id        integer,
    liquid          integer,
    resource_unit   text
);

CREATE TABLE recipes(
    id                    integer PRIMARY KEY AUTO_INCREMENT,
    name                  text,
    users_id              integer,
    production_time_s     integer,
    default_choice        integer
);

CREATE TABLE recipes_inputs(
    id                    integer PRIMARY KEY AUTO_INCREMENT,
    users_id              integer,
    recipes_id            integer,
    resources_id          integer,
    amount                integer,
    FOREIGN KEY(recipes_id) REFERENCES recipes(id)
    ON UPDATE CASCADE ON DELETE SET NULL,
    FOREIGN KEY(resources_id) REFERENCES resources(id)
    ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE TABLE recipes_outputs(
    id                    integer PRIMARY KEY AUTO_INCREMENT,
    users_id              integer,
    recipes_id            integer,
    resources_id          integer,
    amount                integer,
    FOREIGN KEY(recipes_id) REFERENCES recipes(id)
    ON UPDATE CASCADE ON DELETE SET NULL,
    FOREIGN KEY(resources_id) REFERENCES resources(id)
    ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE TABLE machines_recipes(
    id                    integer PRIMARY KEY AUTO_INCREMENT,
    users_id              integer,
    recipes_id           integer,
    machines_id           integer,
    FOREIGN KEY(recipes_id) REFERENCES recipes(id)
    ON UPDATE CASCADE ON DELETE SET NULL,
    FOREIGN KEY(machines_id) REFERENCES machines(id)
    ON UPDATE CASCADE ON DELETE SET NULL
);