CREATE TABLE IF NOT EXISTS songs (
    id serial primary key unique not null,
    song varchar(255) not null,
    group varchar(255) not null,
    text text not null,
    link varchar(255) not null,
    date varchar(255) not null
);