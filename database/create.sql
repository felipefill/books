create table if not exists books(
  id serial primary key,
  title varchar(255),
  description text,
  isbn varchar(13),
  language varchar(2)
);