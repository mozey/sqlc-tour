drop table if exists authors;
create table public.authors (
  id   bigserial primary key,
  name text      not null,
  bio  text
);
