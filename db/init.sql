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
