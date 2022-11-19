insert into genre (title) values ('Рэп');
insert into genre (title) values ('Рок');
insert into genre (title) values ('Поп');
insert into genre (title) values ('Альтернатива');
insert into genre (title) values ('Русские песни');

update genre set human_title = lower(title);

insert into genre_music (music_id, genre_id)
select m.id, 1 as my_genre_id from music as m 
join artist_music as am on am.music_id = m.id
join artist as a on a.id = am.artist_id
where a.artist = 'Post Malone' or
a.artist = 'Eminem' or
a.artist = 'Kendrick Lamar' or
a.artist = 'ЕГОР КРИД';

insert into genre_music (music_id, genre_id)
select m.id, 2 as my_genre_id from music as m 
join artist_music as am on am.music_id = m.id
join artist as a on a.id = am.artist_id
where a.artist = 'Green Day' or
a.artist = 'Queen' or
a.artist = 'Nirvana' or
a.artist = 'Linkin Park' or
a.artist = 'Radiohead' or
a.artist = 'Metallica';

insert into genre_music (music_id, genre_id)
select m.id, 3 as my_genre_id from music as m 
join artist_music as am on am.music_id = m.id
join artist as a on a.id = am.artist_id
where a.artist = 'Taylor Swift' or
a.artist = 'The Weeknd' or
a.artist = 'Zivert' or
a.artist = 'Katy Perry' or
a.artist = 'The Police' or
a.artist = 'Lana Del Rey' or
a.artist = 'ANNA ASTI' or 
a.artist = 'Иван Дорн' or
a.artist = 'Артур Пирожков';

insert into genre_music (music_id, genre_id)
select m.id, 4 as my_genre_id from music as m 
join artist_music as am on am.music_id = m.id
join artist as a on a.id = am.artist_id
where a.artist = 'Imagine Dragons' or
a.artist = 'twenty one pilots' or
a.artist = 'Panic! At The Disco';

insert into genre_music (music_id, genre_id)
select m.id, 5 as my_genre_id from music as m 
join artist_music as am on am.music_id = m.id
join artist as a on a.id = am.artist_id
where a.artist = 'ANNA ASTI' or
a.artist = 'Кино' or
a.artist = 'Григорий Лепс' or
a.artist = 'Земфира' or
a.artist = 'Артур Пирожков' or
a.artist = 'Тима Белорусских' or
a.artist = 'Жуки' or 
a.artist = 'Иван Дорн' or
a.artist = 'Хлеб' or 
a.artist = 'ЕГОР КРИД'or 
a.artist = 'Монеточка';