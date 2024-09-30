CREATE TABLE IF NOT EXISTS songs (
    id serial primary key unique not null,
    song varchar(255) not null,
    group_name varchar(255) not null,
    text text not null,
    link varchar(255) not null,
    date varchar(255) not null
);

CREATE INDEX ix_songs_song ON songs(song);
CREATE INDEX ix_songs_group ON songs(group_name);
CREATE INDEX ix_songs_date ON songs(date);