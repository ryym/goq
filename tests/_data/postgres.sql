-- Do not define constraints such as foreign key
-- to simplify the creation of mock data.

CREATE TABLE IF NOT EXISTS countries (
    id serial NOT NULL,
    name varchar(50) NOT NULL,
    updated_at timestamp without time zone DEFAULT '2000-01-01 09:00:00' NOT NULL
);

CREATE TABLE IF NOT EXISTS cities (
    id serial NOT NULL,
    name varchar(50) NOT NULL,
    country_id smallint NOT NULL,
    updated_at timestamp without time zone DEFAULT '2000-01-01 09:00:00' NOT NULL
);

CREATE TABLE IF NOT EXISTS addresses (
    id serial NOT NULL,
    name varchar(50) NOT NULL,
    city_id smallint NOT NULL,
    updated_at timestamp without time zone DEFAULT '2000-01-01 09:00:00' NOT NULL
);
