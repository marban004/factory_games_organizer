DROP TABLE IF EXISTS users;

CREATE TABLE users(
    id             integer PRIMARY KEY AUTO_INCREMENT,
    login          text,
    passwdhash     text
);