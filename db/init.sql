CREATE TABLE IF NOT EXISTS genre (
    id serial not null unique,
    title varchar not null unique,
    human_title text
);

CREATE TABLE IF NOT EXISTS music (
    id serial not null unique,
    title text,
    artist text,
    duration_two_url text,
    duration_three_url text,
    duration_five_url text,
    duration_fifteen_url text,
    UNIQUE (title, artist)
);

CREATE TABLE IF NOT EXISTS genre_music (
    id serial not null unique,
    genre_id int references "genre"(id) on delete cascade not null,
    music_id int references "music"(id) on delete cascade not null,
    UNIQUE(genre_id, music_id)
);


CREATE TABLE IF NOT EXISTS artist (
    id serial not null unique,
    music_id int references "music"(id) on delete cascade not null,
    artist text unique
);

CREATE TABLE IF NOT EXISTS artist_music (
    id serial not null unique,
    artist_id int references "artist"(id) on delete cascade not null,
    music_id int references "music"(id) on delete cascade not null,
    UNIQUE(artist_id, music_id)
);

CREATE TABLE IF NOT EXISTS human_artist (
    id serial not null unique,
    artist_id int references "artist"(id) on delete cascade not null,
    human_artist text
);

CREATE TABLE IF NOT EXISTS human_title (
    id serial not null unique,
    music_id int references "music"(id) on delete cascade not null,
    human_title text
);

CREATE TABLE IF NOT EXISTS player (
    id serial not null unique,
    vk_id char(64) not null unique,
    points int default 0,
    guessed_songs_count int default 0,
    failed_songs_count int default 0
);

CREATE TABLE IF NOT EXISTS track_history (
    id serial not null unique,
    user_id int references "player"(id),
    track int references "music"(id),
    guessed boolean,
    attempts int
);

DROP SCHEMA public CASCADE;
CREATE SCHEMA public;