create table prefectures (
  id int,
  name varchar(100)
);

create table cities (
  id int,
  name varchar(100),
  prefecture_id int
);

create table towns (
  id int,
  name varchar(100),
  city_id int
);
