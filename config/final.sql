insert into genre (genre) values ('Рэп');
insert into genre (genre) values ('Рок');
insert into genre (genre) values ('Поп');
insert into genre (genre) values ('Альтернатива');
insert into genre (genre) values ('Русские песни');
insert into genre (genre) values ('Зарубежные песни');

update genre set human_genres = array_prepend(lower(genre), human_genres);

insert into genre_music (music_id, genre_id)
select m.id, (select id from genre where genre = 'Рэп') as my_genre_id from music as m 
join artist_music as am on am.music_id = m.id
join artist as a on a.id = am.artist_id
where a.artist_name = 'Post Malone' or
a.artist_name = 'Eminem';

insert into genre_music (music_id, genre_id)
select m.id, (select id from genre where genre = 'Рок') as my_genre_id from music as m 
join artist_music as am on am.music_id = m.id
join artist as a on a.id = am.artist_id
where a.artist_name = 'Green Day' or
a.artist_name = 'Queen' or
a.artist_name = 'Linkin Park' or
a.artist_name = 'Земфира' or
a.artist_name = 'Би-2' or
a.artist_name = 'Кино';

insert into genre_music (music_id, genre_id)
select m.id, (select id from genre where genre = 'Поп') as my_genre_id from music as m 
join artist_music as am on am.music_id = m.id
join artist as a on a.id = am.artist_id
where a.artist_name = 'ЕГОР КРИД' or
a.artist_name = 'The Weeknd' or
a.artist_name = 'Adele' or
a.artist_name = 'Katy Perry' or
a.artist_name = 'ABBA' or
a.artist_name = 'Филипп Киркоров' or
a.artist_name = 'Григорий Лепс' or 
a.artist_name = 'Дискотека Авария' or
a.artist_name = 'Звери' or
a.artist_name = 'Ed Sheeran' or
a.artist_name = 'Руки Вверх';

insert into genre_music (music_id, genre_id)
select m.id, (select id from genre where genre = 'Альтернатива') as my_genre_id from music as m 
join artist_music as am on am.music_id = m.id
join artist as a on a.id = am.artist_id
where a.artist_name = 'Imagine Dragons' or
a.artist_name = 'twenty one pilots';

insert into genre_music (music_id, genre_id)
select m.id, (select id from genre where genre = 'Русские песни') as my_genre_id from music as m 
join artist_music as am on am.music_id = m.id
join artist as a on a.id = am.artist_id
where a.artist_name = 'ЕГОР КРИД' or
a.artist_name = 'Руки Вверх' or
a.artist_name = 'Григорий Лепс' or
a.artist_name = 'Земфира' or
a.artist_name = 'Филипп Киркоров' or
a.artist_name = 'Би-2' or 
a.artist_name = 'Дискотека Авария' or
a.artist_name = 'Кино' or 
a.artist_name = 'Звери';

insert into genre_music (music_id, genre_id)
select m.id, (select id from genre where genre = 'Зарубежные песни') as my_genre_id from music as m 
join artist_music as am on am.music_id = m.id
join artist as a on a.id = am.artist_id
where a.artist_name = 'Post Malone' or
a.artist_name = 'Ed Sheeran' or
a.artist_name = 'Imagine Dragons' or
a.artist_name = 'Queen' or
a.artist_name = 'twenty one pilots' or
a.artist_name = 'Linkin Park' or 
a.artist_name = 'Eminem' or
a.artist_name = 'Green Day' or 
a.artist_name = 'Katy Perry' or
a.artist_name = 'The Weeknd' or
a.artist_name = 'Adele' or 
a.artist_name = 'ABBA';