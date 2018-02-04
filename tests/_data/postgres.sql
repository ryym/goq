-- Do not define constraints such as foreign key
-- to simplify the creation of mock data.

CREATE TABLE IF NOT EXISTS countries (
    id serial NOT NULL,
    name varchar(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS cities (
    id serial NOT NULL,
    name varchar(50) NOT NULL,
    country_id smallint NOT NULL
);

CREATE TABLE IF NOT EXISTS addresses (
    id serial NOT NULL,
    name varchar(50) NOT NULL,
    city_id smallint NOT NULL
);

CREATE TABLE IF NOT EXISTS technologies (
    id serial NOT NULL,
    name varchar(50) NOT NULL,
    description varchar(200) NOT NULL DEFAULT ''
);
