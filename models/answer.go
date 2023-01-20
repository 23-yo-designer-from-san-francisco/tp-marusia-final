package models

// Команды от пользователя
const (
	DontUnderstand  = "не понял"
	AgainE          = "еще раз"
	Again           = "ещё раз"
	ChooseArtist    = "выбрать артиста"
	ChangeArtist    = "поменять артиста"
	ChangeArtist_   = "сменить артиста"
	AnotherArtist   = "другой артист"
	ChooseArtist__  = "выбрать исполнителя"
	ChangeArtist__  = "поменять исполнителя"
	ChangeArtist___ = "сменить исполнителя"
	AnotherArtist__ = "другой исполнитель"
	ChangeGenre     = "поменять жанр"
	ChangeGenre_    = "сменить жанр"
	AnotherGenre    = "другой жанр"
	Genre           = "жанр"
	ChangeGame      = "поменять игру"
	ChangeGame_     = "сменить игру"
	AnotherGame     = "другой игру"
	Next            = "следующ"
	GiveUp          = "сдаюсь"
	Competition     = "соревновани"
	Rule            = "правил"
	Yes             = "да"
	No              = "нет"
	LetsPlay        = "игра"
	Read            = "прочитай"
	Artist          = "исполнител"
	KeyPhrase       = "ключевая фраза"
	IWant           = "Хочу"
)

// Для выбора жанра
const (
	ChooseGenre     = `Назовите жанр, исполнителя или скажите «Любой».` // TODO выбор исполнителя тамже, где и жанр
	List            = "перечисли"
	Available       = "доступные"
	LetsGo          = "давай"
	AvailableGenres = "Я не нашла этот жанр. Вы можете сказать «Любой»."
)

// Правила
const (
	ToMainGame                = "Для одиночной игры скажите «Играем»."
	ToCompetitionGame         = "Чтобы включить соревновательный режим, скажите «Соревнование»."
	ToPhraseGame              = "Чтобы сыграть плейлист по фразе, скажите «Ключевая фраза»."
	PhraseCompetitionRulesTTS = "Для создания плейлиста по фразе, перейдите в миниапп «Угадай мелодию companion». Если угадаете трек с первого раза, дадим 12 баллов. Со второго — 8 баллов. С третьего — 4 балла. Если угадали только исполнителя или название, делим баллы пополам. Играем?"
	PhraseCompetitionRules    = "Для создания плейлиста по фразе, перейдите в миниапп https://vk.com/services?w=app51472048_323199375. Если угадаете трек с первого раза, дадим 12 баллов. Со второго — 8 баллов. С третьего — 4 балла. Если угадали только исполнителя или название, делим баллы пополам. Играем?"
	CompetitionRules          = "После выбора игры, я сделаю соревновательный плейлист. Если угадаете трек с первого раза, дадим 12 баллов. Со второго — 8 баллов. С третьего — 4 балла. Если угадали только исполнителя или название, делим баллы пополам. Играем?"
	ToChange                  = "В любой момент вы можете «сменить игру», «сменить жанр» или «сменить исполнителя»."
	ToStop                    = "Чтобы выйти из навыка, скажите «Стоп»."
)

// Для не тех фраз в определённых местах
const (
	UnknownCommand    = "Я вас не совсем поняла."
	UnknownCommandTTS = "Я вас «НЕ СОВСЕМ» поняла."
	ErrorHappend      = "Что-то мне не хорошо. Попробуйте позже"
)

// Фразы для помощи пользователем с управлением
const (
	Hello             = "Приветствую вас на игре «Угадай мелодию». "
	ToStart           = "Чтобы начать, скажите «Играем»."
	ToRules           = "Чтобы узнать про режимы игры, скажите «Правила»"
	ToContinue        = "Чтобы продолжить игру, скажите «Продолжить»."
	GoodBye           = "Пока. Отлично поиграли!"
	PlaylistFinished  = "Вы сыграли все песни из данного плейлиста."
	KeyPhrasePlaylist = "Назовите ключевую фразу плейлиста."
	KeyPhraseHelp     = "Чтобы узнать подробнее о данном режиме, скажите «Правила»."
)

// Для угадываний
const (
	YourScore = "Ваш счет"
)

var YouDidntGuessTexts = []string{
	"Вы не угадали. ",
	"Неправильно. ",
	"Повезёт в другой раз. ",
}

var YouGuessedTexts = []string{
	"Вы молодец! Угадали!",
	"Абсолютно верно!",
	"Правильно!",
}

var LetsListenTrack = []string{
	"А пока наслаждаемcя песней.",
	"Слушаем полный отрывок.",
	"Можете подпевать!",
}

// Для ответов по песням
const (
	ThatIs = "Это:"
)
