insert into genre (title) values ('Зарубежный реп');
insert into genre (title) values ('Зарубежный рок');
insert into genre (title) values ('Зарубежный поп');

insert into genre_music (music_id, genre_id)
select music.id, 1 as my_genre_id from music join artist on artist.music_id = music.id where artist.human_artist = 'post malone';

insert into genre_music (music_id, genre_id)
select music.id, 1 as my_genre_id from music join artist on artist.music_id = music.id where artist.human_artist = 'kendrick lamar';

insert into genre_music (music_id, genre_id)
select music.id, 1 as my_genre_id from music join artist on artist.music_id = music.id where artist.human_artist = 'eminem';

insert into genre_music (music_id, genre_id)
select music.id, 1 as my_genre_id from music join artist on artist.music_id = music.id where artist.human_artist = 'drake';

insert into genre_music (music_id, genre_id)
select music.id, 2 as my_genre_id from music join artist on artist.music_id = music.id where artist.human_artist = 'ac dc';

insert into genre_music (music_id, genre_id)
select music.id, 2 as my_genre_id from music join artist on artist.music_id = music.id where artist.human_artist = 'rammstein';

insert into genre_music (music_id, genre_id)
select music.id, 2 as my_genre_id from music join artist on artist.music_id = music.id where artist.human_artist = 'green day';

insert into genre_music (music_id, genre_id)
select music.id, 2 as my_genre_id from music join artist on artist.music_id = music.id where artist.human_artist = 'my chemical romance';

insert into genre_music (music_id, genre_id)
select music.id, 2 as my_genre_id from music join artist on artist.music_id = music.id where artist.human_artist = 'nirvana';

insert into genre_music (music_id, genre_id)
select music.id, 3 as my_genre_id from music join artist on artist.music_id = music.id where artist.human_artist = 'sia';

insert into genre_music (music_id, genre_id)
select music.id, 3 as my_genre_id from music join artist on artist.music_id = music.id where artist.human_artist = 'ed sheeran';

insert into genre_music (music_id, genre_id)
select music.id, 3 as my_genre_id from music join artist on artist.music_id = music.id where artist.human_artist = 'rihanna';

insert into genre_music (music_id, genre_id)
select music.id, 3 as my_genre_id from music join artist on artist.music_id = music.id where artist.human_artist = 'billie eilish';

insert into genre_music (music_id, genre_id)
select music.id, 3 as my_genre_id from music join artist on artist.music_id = music.id where artist.human_artist = 'ariana grande';

insert into genre_music (music_id, genre_id)
select music.id, 3 as my_genre_id from music join artist on artist.music_id = music.id where artist.human_artist = 'lana del rey';

update genre set human_title = Lower(title);



