-- +migrate Up
CREATE TABLE users (
    id                       uuid PRIMARY KEY,
    name                     varchar(100) NOT NULL,
    surname                  varchar(100) NOT NULL,
    patronymic               text NOT NULL
);

CREATE TYPE e_pet_enum AS ENUM (
    'dog',
    'cat',
    'other'
    );

CREATE TABLE pets (
    id                       uuid PRIMARY KEY,
    name                     varchar(100) NOT NULL,
    user_id                  uuid NOT NULL,
    pet_type                 e_pet_enum NOT NULL,
    CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS  pets;
DROP TYPE IF EXISTS e_pet_enum;