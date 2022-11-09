CREATE TABLE IF NOT EXISTS music (
    id serial not null unique,
    title text,
    artist text,
    duration_two_url text,
    duration_three_url text,
    duration_five_url text,
    duration_fifteen_url text,
    human_title text,
    UNIQUE (title, artist)
);

CREATE TABLE IF NOT EXISTS artist (
    id serial not null unique,
    music_id int references "music"(id) on delete cascade not null,
    artist text,
    human_artist text
);

CREATE TABLE IF NOT EXISTS users (
    id serial not null unique,
    vk_id char(64) not null unique,
    points int default 0,
    guessed_songs_count int default 0,
    failed_songs_count int default 0
);

CREATE TABLE IF NOT EXISTS track_history (
    id serial not null unique,
    user_id int references "users"(id),
    track int references "music"(id)
);
