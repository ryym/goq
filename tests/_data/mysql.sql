CREATE TABLE IF NOT EXISTS countries (
    id integer primary key AUTO_INCREMENT NOT NULL,
    name varchar(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS cities (
    id integer primary key AUTO_INCREMENT NOT NULL,
    name varchar(50) NOT NULL,
    country_id smallint NOT NULL
);

CREATE TABLE IF NOT EXISTS addresses (
    id integer primary key AUTO_INCREMENT NOT NULL,
    name varchar(50) NOT NULL,
    city_id smallint NOT NULL
);

CREATE TABLE IF NOT EXISTS technologies (
    id integer primary key AUTO_INCREMENT NOT NULL,
    name varchar(50) NOT NULL,
    description varchar(200) NOT NULL DEFAULT ''
);

