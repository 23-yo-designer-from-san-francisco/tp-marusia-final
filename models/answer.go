package models

// Команды от пользователя
const (
	Continue        = "продолжить" // TODO нужно ли? сейчас итак продолжает
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
	CompetitionRule = "правил"
	Yes             = "да"
	No              = "нет"
	LetsPlay        = "игра"
	Read            = "прочитай"
	Artist          = "исполнител"
	RandomPlaylist  = "случайн"
	KeyPhrase       = "ключевая фраза"
	IWant           = "Хочу"
)

// Для выбора жанра
const (
	Notify          = "Вы можете «сменить» игру, жанр или исполнителя - только скажите!"
	ChooseGenre     = `Вы можете выбрать жанр музыки. Назовите жанр или скажите «Любой». Я могу перечислить доступные — только скажите! Если вы хотите выбрать исполнителя, скажите «Исполнители». Чтобы сменить режим игры, скажите «Сменить игру».`
	ChooseGenreTTS  = `Вы можете выбрать ^жанр^ музыки. Назовите жанр или скажите «Любой». Я могу перечислить доступные — только скажите! Если вы хотите выбрать исполнителя, скажите «Исполнители». Чтобы сменить режим игры, скажите «Сменить игру».`
	List            = "перечисли"
	Available       = "доступные"
	LetsGo          = "давай"
	AvailableGenres = "Я не нашла этот жанр. Вы можете сказать «Любой»."
)

// Для соревновательного режима
const (
	CompetitionRules = "После выбора игры, я сделаю соревновательный плейлист. Если угадаете трек с первого раза, дадим 12 баллов. Со второго — 8 баллов. С третьего — 4 балла. Если угадали только исполнителя или название, делим баллы пополам. Играем?"
)

// Для не тех фраз в определённых местах
const (
	AlreadyPlaying    = "Вы уже играете, забыли?" // TODO надо подумать, как это подружить с новыми состояниями
	UnknownCommand    = "Я вас не совсем поняла."
	UnknownCommandTTS = "Я вас «НЕ СОВСЕМ» поняла."
	ErrorHappend      = "Что-то мне не хорошо. Попробуйте позже"
)

// Фразы для помощи пользователем с управлением
const (
	Hello              = "Приветствую вас на игре «Угадай мелодию». "
	ToStart            = "Чтобы начать, скажите «Играем»."
	ToStartCompetitive = "Чтобы начать соревнование, скажите «Соревнование»."
	ToStop             = "Чтобы закончить игру, скажите «Стоп»."
	ToContinue         = "Чтобы продолжить игру, скажите «Продолжить»."
	ToKeyPhrase        = "Чтобы сыграть плейлист по фразе, скажите «Ключевая фраза»"
	GoodBye            = "Пока. Отлично поиграли!"
	PlaylistFinished   = "Вы сыграли все песни из данного плейлиста."
	KeyPhrasePlaylist  = "Назовите ключевую фразу плейлиста."
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
	"Конечно!",
}

// Для ответов по песням
const (
	ThatIs = "Это:"
)
