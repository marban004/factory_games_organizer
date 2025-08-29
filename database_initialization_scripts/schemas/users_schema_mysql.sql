DROP DATABASE IF EXISTS users;
CREATE DATABASE users;
USE users;
DROP TABLE IF EXISTS users;

CREATE TABLE users(
    id             integer PRIMARY KEY AUTO_INCREMENT,
    login          VARCHAR(64),
    passwdhash     text,
    UNIQUE (login)
);