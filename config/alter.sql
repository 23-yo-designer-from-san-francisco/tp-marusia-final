ALTER TABLE artist
  RENAME COLUMN artist TO artist_name;

ALTER TABLE genre
  RENAME COLUMN title TO genre;

ALTER TABLE genre
  RENAME COLUMN human_title TO human_genres;

ALTER TABLE genre
  DROP COLUMN human_genres;

ALTER TABLE genre
  ADD human_genres text[];

UPDATE genre
SET human_genres = array_append(human_genres, lower(genre));