drop table posts cascade if exists;

create table posts (
    id      serial primary key,
    content text,
    author  varchar(255)
);