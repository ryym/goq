-- Do not define constraints such as foreign key
-- to simplify the creation of mock data.

CREATE TABLE IF NOT EXISTS countries (
    id serial NOT NULL,
    name varchar(50) NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS cities (
    id serial NOT NULL,
    name varchar(50) NOT NULL,
    country_id smallint NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);
