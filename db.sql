--
-- PostgreSQL database dump
--

-- Dumped from database version 14.6 (Debian 14.6-1.pgdg110+1)
-- Dumped by pg_dump version 14.6 (Ubuntu 14.6-1.pgdg20.04+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: artist; Type: TABLE; Schema: public; Owner: a_shirshov
--

CREATE TABLE public.artist (
    id integer NOT NULL,
    music_id integer NOT NULL,
    artist text,
    human_artist text
);


ALTER TABLE public.artist OWNER TO a_shirshov;

--
-- Name: artist_id_seq; Type: SEQUENCE; Schema: public; Owner: a_shirshov
--

CREATE SEQUENCE public.artist_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.artist_id_seq OWNER TO a_shirshov;

--
-- Name: artist_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: a_shirshov
--

ALTER SEQUENCE public.artist_id_seq OWNED BY public.artist.id;


--
-- Name: genre; Type: TABLE; Schema: public; Owner: a_shirshov
--

CREATE TABLE public.genre (
    id integer NOT NULL,
    title character varying NOT NULL,
    human_title text
);


ALTER TABLE public.genre OWNER TO a_shirshov;

--
-- Name: genre_id_seq; Type: SEQUENCE; Schema: public; Owner: a_shirshov
--

CREATE SEQUENCE public.genre_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.genre_id_seq OWNER TO a_shirshov;

--
-- Name: genre_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: a_shirshov
--

ALTER SEQUENCE public.genre_id_seq OWNED BY public.genre.id;


--
-- Name: genre_music; Type: TABLE; Schema: public; Owner: a_shirshov
--

CREATE TABLE public.genre_music (
    id integer NOT NULL,
    genre_id integer NOT NULL,
    music_id integer NOT NULL
);


ALTER TABLE public.genre_music OWNER TO a_shirshov;

--
-- Name: genre_music_id_seq; Type: SEQUENCE; Schema: public; Owner: a_shirshov
--

CREATE SEQUENCE public.genre_music_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.genre_music_id_seq OWNER TO a_shirshov;

--
-- Name: genre_music_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: a_shirshov
--

ALTER SEQUENCE public.genre_music_id_seq OWNED BY public.genre_music.id;


--
-- Name: music; Type: TABLE; Schema: public; Owner: a_shirshov
--

CREATE TABLE public.music (
    id integer NOT NULL,
    title text,
    artist text,
    duration_two_url text,
    duration_three_url text,
    duration_five_url text,
    duration_fifteen_url text,
    human_title text
);


ALTER TABLE public.music OWNER TO a_shirshov;

--
-- Name: music_id_seq; Type: SEQUENCE; Schema: public; Owner: a_shirshov
--

CREATE SEQUENCE public.music_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.music_id_seq OWNER TO a_shirshov;

--
-- Name: music_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: a_shirshov
--

ALTER SEQUENCE public.music_id_seq OWNED BY public.music.id;


--
-- Name: player; Type: TABLE; Schema: public; Owner: a_shirshov
--

CREATE TABLE public.player (
    id integer NOT NULL,
    vk_id character(64) NOT NULL,
    points integer DEFAULT 0,
    guessed_songs_count integer DEFAULT 0,
    failed_songs_count integer DEFAULT 0
);


ALTER TABLE public.player OWNER TO a_shirshov;

--
-- Name: player_id_seq; Type: SEQUENCE; Schema: public; Owner: a_shirshov
--

CREATE SEQUENCE public.player_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.player_id_seq OWNER TO a_shirshov;

--
-- Name: player_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: a_shirshov
--

ALTER SEQUENCE public.player_id_seq OWNED BY public.player.id;


--
-- Name: track_history; Type: TABLE; Schema: public; Owner: a_shirshov
--

CREATE TABLE public.track_history (
    id integer NOT NULL,
    user_id integer,
    track integer,
    guessed boolean,
    attempts integer
);


ALTER TABLE public.track_history OWNER TO a_shirshov;

--
-- Name: track_history_id_seq; Type: SEQUENCE; Schema: public; Owner: a_shirshov
--

CREATE SEQUENCE public.track_history_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.track_history_id_seq OWNER TO a_shirshov;

--
-- Name: track_history_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: a_shirshov
--

ALTER SEQUENCE public.track_history_id_seq OWNED BY public.track_history.id;


--
-- Name: artist id; Type: DEFAULT; Schema: public; Owner: a_shirshov
--

ALTER TABLE ONLY public.artist ALTER COLUMN id SET DEFAULT nextval('public.artist_id_seq'::regclass);


--
-- Name: genre id; Type: DEFAULT; Schema: public; Owner: a_shirshov
--

ALTER TABLE ONLY public.genre ALTER COLUMN id SET DEFAULT nextval('public.genre_id_seq'::regclass);


--
-- Name: genre_music id; Type: DEFAULT; Schema: public; Owner: a_shirshov
--

ALTER TABLE ONLY public.genre_music ALTER COLUMN id SET DEFAULT nextval('public.genre_music_id_seq'::regclass);


--
-- Name: music id; Type: DEFAULT; Schema: public; Owner: a_shirshov
--

ALTER TABLE ONLY public.music ALTER COLUMN id SET DEFAULT nextval('public.music_id_seq'::regclass);


--
-- Name: player id; Type: DEFAULT; Schema: public; Owner: a_shirshov
--

ALTER TABLE ONLY public.player ALTER COLUMN id SET DEFAULT nextval('public.player_id_seq'::regclass);


--
-- Name: track_history id; Type: DEFAULT; Schema: public; Owner: a_shirshov
--

ALTER TABLE ONLY public.track_history ALTER COLUMN id SET DEFAULT nextval('public.track_history_id_seq'::regclass);


--
-- Data for Name: artist; Type: TABLE DATA; Schema: public; Owner: a_shirshov
--

COPY public.artist (id, music_id, artist, human_artist) FROM stdin;
1	1	Kendrick Lamar	kendrick lamar
2	1	SZA	sza
3	2	Kendrick Lamar	kendrick lamar
4	3	Green Day	green day
5	4	Queen	queen
6	5	Amy Winehouse	amy winehouse
7	6	Green Day	green day
8	7	Imagine Dragons	imagine dragons
9	8	Post Malone	post malone
10	9	Taylor Swift	taylor swift
11	10	The Weeknd	the weeknd
12	11	Imagine Dragons	imagine dragons
13	12	Green Day	green day
14	13	The Weeknd	the weeknd
15	14	Post Malone	post malone
16	15	Nirvana	nirvana
17	16	Linkin Park	linkin park
18	17	Zivert	zivert
19	18	Radiohead	radiohead
20	19	Fall Out Boy	fall out boy
21	20	Katy Perry	katy perry
22	20	Juicy J	juicy j
23	21	Imagine Dragons	imagine dragons
24	21	JID	jid
25	21	Arcane	arcane
26	21	League Of Legends	league of legends
27	22	Metallica	metallica
28	23	The Police	the police
29	24	Backstreet Boys	backstreet boys
30	25	Linkin Park	linkin park
31	26	twenty one pilots	twenty one pilots
32	27	Katy Perry	katy perry
33	28	Metallica	metallica
34	29	Daft Punk	daft punk
35	30	Kendrick Lamar	kendrick lamar
36	31	Panic! At The Disco	panic at the disco
37	32	Green Day	green day
38	33	twenty one pilots	twenty one pilots
39	34	Panic! At The Disco	panic at the disco
40	35	Miyagi & Эндшпиль	miyagi & эндшпиль
41	35	Рем Дигга	рем дигга
42	36	Katy Perry	katy perry
43	37	Taylor Swift	taylor swift
44	38	Whitney Houston	whitney houston
45	39	Queen	queen
46	40	Panic! At The Disco	panic at the disco
47	41	Nirvana	nirvana
48	42	The Weeknd	the weeknd
49	43	Linkin Park	linkin park
50	44	Bon Jovi	bon jovi
51	45	Queen	queen
52	46	Kendrick Lamar	kendrick lamar
53	46	Zacari	zacari
54	47	Katy Perry	katy perry
55	48	Zivert	zivert
56	49	System Of A Down	system of a down
57	50	Eminem	eminem
58	51	Queen	queen
59	52	Taylor Swift	taylor swift
60	53	Metallica	metallica
61	54	Eminem	eminem
62	55	Kendrick Lamar	kendrick lamar
63	55	Jay Rock	jay rock
64	56	Skillet	skillet
65	57	twenty one pilots	twenty one pilots
66	58	Imagine Dragons	imagine dragons
67	59	Linkin Park	linkin park
68	60	Panic! At The Disco	panic at the disco
69	61	Metallica	metallica
70	62	Linkin Park	linkin park
71	63	Arctic Monkeys	arctic monkeys
72	64	Eminem	eminem
73	65	Nirvana	nirvana
74	66	Twenty One Pilots	twenty one pilots
75	67	Katy Perry	katy perry
76	68	The Weeknd	the weeknd
77	69	Nirvana	nirvana
78	70	Queen	queen
79	71	The Killers	the killers
80	72	Nirvana	nirvana
81	73	The Weeknd	the weeknd
82	73	Daft Punk	daft punk
83	74	Post Malone	post malone
84	75	Lana Del Rey	lana del rey
85	76	Eminem	eminem
86	76	Dina Rae	dina rae
87	77	Hozier	hozier
88	78	twenty one pilots	twenty one pilots
89	79	Metallica	metallica
90	80	The Prodigy	the prodigy
91	81	Green Day	green day
92	82	Imagine Dragons	imagine dragons
93	83	Taylor Swift	taylor swift
94	84	Eminem	eminem
95	85	Post Malone	post malone
96	86	Taylor Swift	taylor swift
97	87	Lana Del Rey	lana del rey
98	88	Post Malone	post malone
99	88	21 Savage	21 savage
100	89	Жуки	жуки
101	90	Тима Белорусских	тима белорусских
102	91	Митя Фомин	митя фомин
103	92	Звери	звери
104	93	Miyagi	miyagi
105	93	Andy Panda	andy panda
106	94	Кино	кино
107	94	В. Цой	в. цой
108	95	Zivert	zivert
109	95	M'Dee	m'dee
110	96	Звери	звери
111	97	Звери	звери
112	98	Дора	дора
113	99	Земфира	земфира
114	100	Артур Пирожков	артур пирожков
115	101	В. Цой	в. цой
116	101	Кино	кино
117	102	Zivert	zivert
118	103	Григорий Лепс	григорий лепс
119	103	Ани Лорак	ани лорак
120	104	МакSим	макsим
121	105	Земфира	земфира
122	106	Монеточка	монеточка
123	107	Иван Дорн	иван дорн
124	108	Григорий Лепс	григорий лепс
125	108	Тимати	тимати
126	109	Ленинград	ленинград
127	110	Тима Белорусских	тима белорусских
128	111	Валентин Стрыкало	валентин стрыкало
129	112	ANNA ASTI	anna asti
130	113	Miyagi & Andy Panda	miyagi & andy panda
131	114	В. Цой	в. цой
132	114	Кино	кино
133	115	Кино	кино
134	115	В. Цой	в. цой
135	116	ANNA ASTI	anna asti
136	117	ANNA ASTI	anna asti
137	118	Земфира	земфира
138	119	Звери	звери
139	120	Звери	звери
140	121	Miyagi	miyagi
141	121	KADI	kadi
142	122	Григорий Лепс	григорий лепс
143	123	ЕГОР КРИД	егор крид
144	124	Леша Свик	леша свик
145	125	В. Цой	в. цой
146	125	Кино	кино
147	126	Miyagi & Andy Panda	miyagi & andy panda
148	127	ANNA ASTI	anna asti
149	128	Земфира	земфира
150	129	ANNA ASTI	anna asti
151	129	Филипп Киркоров	филипп киркоров
152	130	Григорий Лепс	григорий лепс
153	131	Хлеб	хлеб
154	132	Григорий Лепс	григорий лепс
155	133	Земфира	земфира
156	134	Zivert	zivert
\.


--
-- Data for Name: genre; Type: TABLE DATA; Schema: public; Owner: a_shirshov
--

COPY public.genre (id, title, human_title) FROM stdin;
1	Реп	\N
2	Рок	\N
3	Поп	\N
4	Альтернатива	\N
5	Русские песни	\N
\.


--
-- Data for Name: genre_music; Type: TABLE DATA; Schema: public; Owner: a_shirshov
--

COPY public.genre_music (id, genre_id, music_id) FROM stdin;
1	1	8
2	1	14
3	1	74
4	1	85
5	1	88
6	1	50
7	1	54
8	1	64
9	1	76
10	1	84
11	1	1
12	1	2
13	1	30
14	1	46
15	1	55
16	2	3
17	2	4
18	2	6
19	2	12
20	2	15
21	2	16
22	2	18
23	2	25
24	2	32
25	2	39
26	2	41
27	2	43
28	2	45
29	2	51
30	2	59
31	2	62
32	2	65
33	2	69
34	2	70
35	2	72
36	2	81
37	3	9
38	3	10
39	3	13
40	3	20
41	3	27
42	3	36
43	3	37
44	3	42
45	3	47
46	3	52
47	3	67
48	3	68
49	3	73
50	3	75
51	3	83
52	3	86
53	3	87
54	4	7
55	4	11
56	4	21
57	4	26
58	4	31
59	4	33
60	4	34
61	4	40
62	4	57
63	4	58
64	4	60
65	4	66
66	4	78
67	4	82
68	5	17
69	5	35
70	5	48
71	5	92
72	5	94
73	5	95
74	5	96
75	5	97
76	5	101
77	5	102
78	5	103
79	5	108
80	5	112
81	5	113
82	5	114
83	5	115
84	5	116
85	5	117
86	5	119
87	5	120
88	5	122
89	5	125
90	5	126
91	5	127
92	5	129
93	5	130
94	5	132
95	5	134
\.


--
-- Data for Name: music; Type: TABLE DATA; Schema: public; Owner: a_shirshov
--

COPY public.music (id, title, artist, duration_two_url, duration_three_url, duration_five_url, duration_fifteen_url, human_title) FROM stdin;
1	All The Stars	Kendrick Lamar, SZA	2000512035_456239220	2000512035_456239219	2000512035_456239223		all the stars
2	Alright	Kendrick Lamar	2000512035_456239207	2000512035_456239208	2000512035_456239209		alright
3	American Idiot	Green Day	2000512035_456239159	2000512035_456239160	2000512035_456239161		american idiot
4	Another One Bites The Dust	Queen	2000512035_456239303	2000512035_456239304	2000512035_456239305		another one bites the dust
5	Back To Black	Amy Winehouse	2000512035_456239108	2000512035_456239109	2000512035_456239110		back to black
6	Basket Case	Green Day	2000512035_456239162	2000512035_456239163	2000512035_456239164		basket case
7	Believer	Imagine Dragons	2000512035_456239177	2000512035_456239178	2000512035_456239179		believer
8	Better Now	Post Malone	2000512035_456239287	2000512035_456239289	2000512035_456239290		better now
9	Blank Space	Taylor Swift	2000512035_456239327	2000512035_456239328	2000512035_456239329		blank space
10	Blinding Lights	The Weeknd	2000512035_456239351	2000512035_456239352	2000512035_456239353		blinding lights
11	Bones	Imagine Dragons	2000512035_456239180	2000512035_456239181	2000512035_456239182		bones
12	Boulevard of Broken Dreams	Green Day	2000512035_456239165	2000512035_456239166	2000512035_456239167		boulevard of broken dreams
13	Call Out My Name	The Weeknd	2000512035_456239354	2000512035_456239355	2000512035_456239356		call out my name
14	Circles	Post Malone	2000512035_456239291	2000512035_456239292	2000512035_456239293		circles
15	Come As You Are	Nirvana	2000512035_456239273	2000512035_456239274	2000512035_456239275		come as you are
16	Crawling	Linkin Park	2000512035_456239228	2000512035_456239229	2000512035_456239230		crawling
17	Credo	Zivert	2000512035_456239384	2000512035_456239385	2000512035_456239386		credo
18	Creep	Radiohead	2000512035_456239317	2000512035_456239319	2000512035_456239320		creep
19	Dance, Dance	Fall Out Boy	2000512035_456239156	2000512035_456239157	2000512035_456239158		dance dance
20	Dark Horse	Katy Perry feat. Juicy J	2000512035_456239204	2000512035_456239206	2000512035_456239205		dark horse
21	Enemy	Imagine Dragons, JID, Arcane, League Of Legends	2000512035_456239189	2000512035_456239190	2000512035_456239191		enemy
22	Enter Sandman	Metallica	2000512035_456239243	2000512035_456239244	2000512035_456239245		enter sandman
23	Every Breath You Take	The Police	2000512035_456239345	2000512035_456239346	2000512035_456239347		every breath you take
24	Everybody 	Backstreet Boys	2000512035_456239129	2000512035_456239130	2000512035_456239131		everybody 
25	Faint	Linkin Park	2000512035_456239232	2000512035_456239231	2000512035_456239234		faint
26	Fake You Out	twenty one pilots	2000512035_456239366	2000512035_456239367	2000512035_456239368		fake you out
27	Firework	Katy Perry	2000512035_456239192	2000512035_456239196	2000512035_456239195		firework
28	Fuel	Metallica	2000512035_456239246	2000512035_456239247	2000512035_456239248		fuel
29	Get Lucky	Daft Punk	2000512035_456239138	2000512035_456239139	2000512035_456239140		get lucky
30	HUMBLE.	Kendrick Lamar	2000512035_456239210	2000512035_456239211	2000512035_456239212		humble
31	High Hopes	Panic! At The Disco	2000512035_456239498	2000512035_456239499	2000512035_456239500		high hopes
32	Holiday	Green Day	2000512035_456239168	2000512035_456239169	2000512035_456239170		holiday
33	Hometown	twenty one pilots	2000512035_456239369	2000512035_456239370	2000512035_456239371		hometown
34	House of Memories	Panic! At The Disco	2000512035_456239501	2000512035_456239502	2000512035_456239503		house of memories
35	I Got Love	Miyagi & Эндшпиль feat. Рем Дигга	2000512035_456239264	2000512035_456239265	2000512035_456239266		i got love
36	I Kissed A Girl	Katy Perry	2000512035_456239193	2000512035_456239194	2000512035_456239197		i kissed a girl
37	I Knew You Were Trouble.	Taylor Swift	2000512035_456239330	2000512035_456239331	2000512035_456239332		i knew you were trouble
38	I Wanna Dance with Somebody 	Whitney Houston	2000512035_456239381	2000512035_456239382	2000512035_456239383		i wanna dance with somebody 
39	I Want To Break Free	Queen	2000512035_456239306	2000512035_456239307	2000512035_456239308		i want to break free
40	I Write Sins Not Tragedies	Panic! At The Disco	2000512035_456239504	2000512035_456239505	2000512035_456239506		i write sins not tragedies
41	In Bloom	Nirvana	2000512035_456239276	2000512035_456239277	2000512035_456239278		in bloom
42	In Your Eyes	The Weeknd	2000512035_456239357	2000512035_456239358	2000512035_456239359		in your eyes
43	In the End	Linkin Park	2000512035_456239233	2000512035_456239235	2000512035_456239236		in the end
44	It's My Life	Bon Jovi	2000512035_456239132	2000512035_456239133	2000512035_456239134		it's my life
45	Killer Queen	Queen	2000512035_456239309	2000512035_456239310	2000512035_456239311		killer queen
46	LOVE.	Kendrick Lamar feat. Zacari	2000512035_456239216	2000512035_456239217	2000512035_456239218		love
47	Last Friday Night 	Katy Perry	2000512035_456239198	2000512035_456239199	2000512035_456239200		last friday night 
48	Life	Zivert	2000512035_456239387	2000512035_456239388	2000512035_456239389		life
49	Lonely Day	System Of A Down	2000512035_456239324	2000512035_456239325	2000512035_456239326		lonely day
50	Lose Yourself	Eminem	2000512035_456239141	2000512035_456239142	2000512035_456239143		lose yourself
51	Love Of My Life	Queen	2000512035_456239312	2000512035_456239313	2000512035_456239314		love of my life
52	Love Story	Taylor Swift	2000512035_456239333	2000512035_456239335	2000512035_456239334		love story
53	Master Of Puppets	Metallica	2000512035_456239249	2000512035_456239250	2000512035_456239251		master of puppets
54	Mockingbird	Eminem	2000512035_456239144	2000512035_456239145	2000512035_456239146		mockingbird
55	Money Trees	Kendrick Lamar feat. Jay Rock	2000512035_456239213	2000512035_456239214	2000512035_456239215		money trees
56	Monster	Skillet	2000512035_456239321	2000512035_456239322	2000512035_456239323		monster
57	My Blood	twenty one pilots	2000512035_456239372	2000512035_456239373	2000512035_456239374		my blood
58	Natural	Imagine Dragons	2000512035_456239183	2000512035_456239184	2000512035_456239185		natural
59	New Divide	Linkin Park	2000512035_456239237	2000512035_456239238	2000512035_456239239		new divide
60	Nicotine	Panic! At The Disco	2000512035_456239507	2000512035_456239508	2000512035_456239509		nicotine
61	Nothing Else Matters	Metallica	2000512035_456239252	2000512035_456239253	2000512035_456239255		nothing else matters
62	Numb	Linkin Park	2000512035_456239240	2000512035_456239241	2000512035_456239242		numb
63	Old Yellow Bricks	Arctic Monkeys	2000512035_456239126	2000512035_456239127	2000512035_456239128		old yellow bricks
64	Rap God	Eminem	2000512035_456239147	2000512035_456239148	2000512035_456239149		rap god
65	Rape Me	Nirvana	2000512035_456239279	2000512035_456239280	2000512035_456239281		rape me
66	Ride	Twenty One Pilots	2000512035_456239375	2000512035_456239376	2000512035_456239377		ride
67	Roar	Katy Perry	2000512035_456239201	2000512035_456239202	2000512035_456239203		roar
68	Save Your Tears	The Weeknd	2000512035_456239360	2000512035_456239361	2000512035_456239362		save your tears
69	Smells Like Teen Spirit	Nirvana	2000512035_456239282	2000512035_456239283	2000512035_456239284		smells like teen spirit
70	Somebody To Love	Queen	2000512035_456239315	2000512035_456239316	2000512035_456239318		somebody to love
71	Somebody Told Me	The Killers	2000512035_456239342	2000512035_456239343	2000512035_456239344		somebody told me
72	Something In The Way	Nirvana	2000512035_456239285	2000512035_456239286	2000512035_456239288		something in the way
73	Starboy	The Weeknd feat. Daft Punk	2000512035_456239363	2000512035_456239364	2000512035_456239365		starboy
74	Stay	Post Malone	2000512035_456239294	2000512035_456239295	2000512035_456239296		stay
75	Summertime Sadness	Lana Del Rey	2000512035_456239221	2000512035_456239224	2000512035_456239222		summertime sadness
76	Superman	Eminem feat. Dina Rae	2000512035_456239153	2000512035_456239154	2000512035_456239155		superman
77	Take Me To Church	Hozier	2000512035_456239173	2000512035_456239174	2000512035_456239176		take me to church
78	Tear in My Heart	twenty one pilots	2000512035_456239378	2000512035_456239379	2000512035_456239380		tear in my heart
79	The Unforgiven	Metallica	2000512035_456239254	2000512035_456239256	2000512035_456239257		the unforgiven
80	Voodoo People	The Prodigy	2000512035_456239348	2000512035_456239349	2000512035_456239350		voodoo people
81	Wake Me up When September Ends	Green Day	2000512035_456239171	2000512035_456239172	2000512035_456239175		wake me up when september ends
82	Whatever It Takes	Imagine Dragons	2000512035_456239186	2000512035_456239187	2000512035_456239188		whatever it takes
83	Wildest Dreams	Taylor Swift	2000512035_456239337	2000512035_456239336	2000512035_456239338		wildest dreams
84	Without Me	Eminem	2000512035_456239150	2000512035_456239151	2000512035_456239152		without me
85	Wow.	Post Malone	2000512035_456239297	2000512035_456239298	2000512035_456239299		wow
86	You Belong With Me	Taylor Swift	2000512035_456239340	2000512035_456239339	2000512035_456239341		you belong with me
87	Young And Beautiful	Lana Del Rey	2000512035_456239225	2000512035_456239227	2000512035_456239226		young and beautiful
88	rockstar	Post Malone feat. 21 Savage	2000512035_456239300	2000512035_456239301	2000512035_456239302		rockstar
89	Батарейка	Жуки	2000512035_456239435	2000512035_456239436	2000512035_456239437		батарейка
90	Витаминка	Тима Белорусских	2000512035_456239483	2000512035_456239484	2000512035_456239485		витаминка
91	Все будет хорошо	Митя Фомин	2000512035_456239477	2000512035_456239478	2000512035_456239479		все будет хорошо
92	Все, что касается	Звери	2000512035_456239438	2000512035_456239439	2000512035_456239440		все что касается
93	Говори мне	Miyagi feat. Andy Panda	2000512035_456239267	2000512035_456239268	2000512035_456239269		говори мне
94	Группа крови	Кино, В. Цой	2000512035_456239462	2000512035_456239463	2000512035_456239464		группа крови
95	Двусмысленно	Zivert feat. M'Dee	2000512035_456239396	2000512035_456239398	2000512035_456239397		двусмысленно
96	До скорой встречи!	Звери	2000512035_456239441	2000512035_456239442	2000512035_456239444		до скорой встречи
97	Дожди-пистолеты	Звери	2000512035_456239443	2000512035_456239445	2000512035_456239446		дожди-пистолеты
98	Дорадура	Дора	2000512035_456239429	2000512035_456239430	2000512035_456239431		дорадура
99	ЖИТЬ В ТВОЕЙ ГОЛОВЕ	Земфира	2000512035_456239495	2000512035_456239496	2000512035_456239497		жить в твоей голове
100	Зацепила	Артур Пирожков	2000512035_456239399	2000512035_456239400	2000512035_456239401		зацепила
101	Звезда по имени Солнце	В. Цой, Кино	2000512035_456239402	2000512035_456239403	2000512035_456239404		звезда по имени солнце
102	Зелёные волны	Zivert	2000512035_456239390	2000512035_456239391	2000512035_456239392		зелёные волны
103	Зеркала	Григорий Лепс, Ани Лорак	2000512035_456239422	2000512035_456239424	2000512035_456239428		зеркала
104	Знаешь ли ты	МакSим	2000512035_456239474	2000512035_456239475	2000512035_456239476		знаешь ли ты
105	ИСКАЛА	Земфира	2000512035_456239135	2000512035_456239136	2000512035_456239137		искала
106	Каждый раз	Монеточка	2000512035_456239480	2000512035_456239481	2000512035_456239482		каждый раз
107	Кричу	Иван Дорн	2000512035_456239459	2000512035_456239460	2000512035_456239461		кричу
108	Лондон	Григорий Лепс, Тимати	2000512035_456239426	2000512035_456239425	2000512035_456239427		лондон
109	Менеджер	Ленинград	2000512035_456239467	2000512035_456239469	2000512035_456239470		менеджер
110	Мокрые кроссы	Тима Белорусских	2000512035_456239486	2000512035_456239487	2000512035_456239489		мокрые кроссы
111	Наше лето	Валентин Стрыкало	2000512035_456239411	2000512035_456239412	2000512035_456239413		наше лето
112	Ночью на кухне	ANNA ASTI	2000512035_456239111	2000512035_456239112	2000512035_456239113		ночью на кухне
113	Патрон	Miyagi & Andy Panda	2000512035_456239258	2000512035_456239259	2000512035_456239260		патрон
114	Пачка сигарет	В. Цой, Кино	2000512035_456239405	2000512035_456239406	2000512035_456239407		пачка сигарет
115	Перемен	Кино, В. Цой	2000512035_456239465	2000512035_456239466	2000512035_456239468		перемен
116	По барам	ANNA ASTI	2000512035_456239114	2000512035_456239115	2000512035_456239116		по барам
117	Повело	ANNA ASTI	2000512035_456239117	2000512035_456239118	2000512035_456239119		повело
118	Прогулка	Земфира	2000512035_456239453	2000512035_456239454	2000512035_456239455		прогулка
119	Просто такая сильная любовь	Звери	2000512035_456239448	2000512035_456239447	2000512035_456239449		просто такая сильная любовь
120	Районы-кварталы	Звери	2000512035_456239450	2000512035_456239451	2000512035_456239452		районы-кварталы
121	Родная пой	Miyagi feat. KADI	2000512035_456239270	2000512035_456239271	2000512035_456239272		родная пой
122	Рюмка водки на столе	Григорий Лепс	2000512035_456239414	2000512035_456239415	2000512035_456239416		рюмка водки на столе
123	Самая самая	ЕГОР КРИД	2000512035_456239432	2000512035_456239433	2000512035_456239434		самая самая
124	Самолеты	Леша Свик	2000512035_456239473	2000512035_456239471	2000512035_456239472		самолеты
125	Спокойная ночь	В. Цой, Кино	2000512035_456239408	2000512035_456239409	2000512035_456239410		спокойная ночь
126	Там ревели горы	Miyagi & Andy Panda	2000512035_456239262	2000512035_456239261	2000512035_456239263		там ревели горы
127	Феникс	ANNA ASTI	2000512035_456239120	2000512035_456239121	2000512035_456239122		феникс
128	ХОЧЕШЬ	Земфира	2000512035_456239492	2000512035_456239493	2000512035_456239494		хочешь
129	Хобби	ANNA ASTI, Филипп Киркоров	2000512035_456239123	2000512035_456239124	2000512035_456239125		хобби
130	Что ж ты натворила	Григорий Лепс	2000512035_456239418	2000512035_456239420	2000512035_456239417		что ж ты натворила
131	Шашлындос	Хлеб	2000512035_456239488	2000512035_456239491	2000512035_456239490		шашлындос
132	Я счастливый	Григорий Лепс	2000512035_456239419	2000512035_456239421	2000512035_456239423		я счастливый
133	ромашки	Земфира	2000512035_456239456	2000512035_456239457	2000512035_456239458		ромашки
134	смех и грех	Zivert	2000512035_456239393	2000512035_456239394	2000512035_456239395		смех и грех
\.


--
-- Data for Name: player; Type: TABLE DATA; Schema: public; Owner: a_shirshov
--

COPY public.player (id, vk_id, points, guessed_songs_count, failed_songs_count) FROM stdin;
\.


--
-- Data for Name: track_history; Type: TABLE DATA; Schema: public; Owner: a_shirshov
--

COPY public.track_history (id, user_id, track, guessed, attempts) FROM stdin;
\.


--
-- Name: artist_id_seq; Type: SEQUENCE SET; Schema: public; Owner: a_shirshov
--

SELECT pg_catalog.setval('public.artist_id_seq', 156, true);


--
-- Name: genre_id_seq; Type: SEQUENCE SET; Schema: public; Owner: a_shirshov
--

SELECT pg_catalog.setval('public.genre_id_seq', 5, true);


--
-- Name: genre_music_id_seq; Type: SEQUENCE SET; Schema: public; Owner: a_shirshov
--

SELECT pg_catalog.setval('public.genre_music_id_seq', 95, true);


--
-- Name: music_id_seq; Type: SEQUENCE SET; Schema: public; Owner: a_shirshov
--

SELECT pg_catalog.setval('public.music_id_seq', 134, true);


--
-- Name: player_id_seq; Type: SEQUENCE SET; Schema: public; Owner: a_shirshov
--

SELECT pg_catalog.setval('public.player_id_seq', 1, false);


--
-- Name: track_history_id_seq; Type: SEQUENCE SET; Schema: public; Owner: a_shirshov
--

SELECT pg_catalog.setval('public.track_history_id_seq', 1, false);


--
-- Name: artist artist_id_key; Type: CONSTRAINT; Schema: public; Owner: a_shirshov
--

ALTER TABLE ONLY public.artist
    ADD CONSTRAINT artist_id_key UNIQUE (id);


--
-- Name: genre genre_id_key; Type: CONSTRAINT; Schema: public; Owner: a_shirshov
--

ALTER TABLE ONLY public.genre
    ADD CONSTRAINT genre_id_key UNIQUE (id);


--
-- Name: genre_music genre_music_genre_id_music_id_key; Type: CONSTRAINT; Schema: public; Owner: a_shirshov
--

ALTER TABLE ONLY public.genre_music
    ADD CONSTRAINT genre_music_genre_id_music_id_key UNIQUE (genre_id, music_id);


--
-- Name: genre_music genre_music_id_key; Type: CONSTRAINT; Schema: public; Owner: a_shirshov
--

ALTER TABLE ONLY public.genre_music
    ADD CONSTRAINT genre_music_id_key UNIQUE (id);


--
-- Name: genre genre_title_key; Type: CONSTRAINT; Schema: public; Owner: a_shirshov
--

ALTER TABLE ONLY public.genre
    ADD CONSTRAINT genre_title_key UNIQUE (title);


--
-- Name: music music_id_key; Type: CONSTRAINT; Schema: public; Owner: a_shirshov
--

ALTER TABLE ONLY public.music
    ADD CONSTRAINT music_id_key UNIQUE (id);


--
-- Name: music music_title_artist_key; Type: CONSTRAINT; Schema: public; Owner: a_shirshov
--

ALTER TABLE ONLY public.music
    ADD CONSTRAINT music_title_artist_key UNIQUE (title, artist);


--
-- Name: player player_id_key; Type: CONSTRAINT; Schema: public; Owner: a_shirshov
--

ALTER TABLE ONLY public.player
    ADD CONSTRAINT player_id_key UNIQUE (id);


--
-- Name: player player_vk_id_key; Type: CONSTRAINT; Schema: public; Owner: a_shirshov
--

ALTER TABLE ONLY public.player
    ADD CONSTRAINT player_vk_id_key UNIQUE (vk_id);


--
-- Name: track_history track_history_id_key; Type: CONSTRAINT; Schema: public; Owner: a_shirshov
--

ALTER TABLE ONLY public.track_history
    ADD CONSTRAINT track_history_id_key UNIQUE (id);


--
-- Name: artist artist_music_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: a_shirshov
--

ALTER TABLE ONLY public.artist
    ADD CONSTRAINT artist_music_id_fkey FOREIGN KEY (music_id) REFERENCES public.music(id) ON DELETE CASCADE;


--
-- Name: genre_music genre_music_genre_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: a_shirshov
--

ALTER TABLE ONLY public.genre_music
    ADD CONSTRAINT genre_music_genre_id_fkey FOREIGN KEY (genre_id) REFERENCES public.genre(id) ON DELETE CASCADE;


--
-- Name: genre_music genre_music_music_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: a_shirshov
--

ALTER TABLE ONLY public.genre_music
    ADD CONSTRAINT genre_music_music_id_fkey FOREIGN KEY (music_id) REFERENCES public.music(id) ON DELETE CASCADE;


--
-- Name: track_history track_history_track_fkey; Type: FK CONSTRAINT; Schema: public; Owner: a_shirshov
--

ALTER TABLE ONLY public.track_history
    ADD CONSTRAINT track_history_track_fkey FOREIGN KEY (track) REFERENCES public.music(id);


--
-- Name: track_history track_history_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: a_shirshov
--

ALTER TABLE ONLY public.track_history
    ADD CONSTRAINT track_history_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.player(id);


--
-- PostgreSQL database dump complete
--

