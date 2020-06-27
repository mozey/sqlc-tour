drop table if exists authors;
create table public.authors
(
    id   bigserial primary key,
    name varchar unique not null,
    bio  text
);
