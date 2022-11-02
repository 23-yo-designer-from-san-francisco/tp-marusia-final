package answer

const (
	Continue       string = "продолжить" // TODO нужно ли? сейчас итак продолжает
	Repeat                = "повтори"    // TODO почему-то не работает
	DontUnderstand        = "не понял"
	AgainE                = "еще раз"
	Again                 = "ещё раз"
	ChangeGenre           = "поменять жанр"
	ChangeGenre_          = "сменить жанр"
	AnotherGenre          = "другой жанр"
)

// Для выбора жанра
const (
	ChooseGenre     string = "Вы можете выбрать жанр музыки. Назовите жанр или скажите «Любой». Я могу перечислить доступные — только скажите!"
	List            string = "перечисли"
	Available       string = "доступн"
	LetsGo          string = "давай"
	AvailableGenres string = "Сейчас доступны жанры: «Рок» и «Не Рок». Либо скажите — «Любой жанр»."
	Rock            string = "Рок"
	NotRock         string = "Не рок"
	Any             string = "Любой"
)

// Для не тех фраз в определённых местах
const (
	AlreadyPlaying    = "Вы уже играете, забыли? " // TODO надо подумать, как это подружить с новыми состояниями
	UnknownCommand    = "Я вас не совсем поняла. "
	UnknownCommandTTS = "Я вас «НЕ СОВСЕМ» поняла. "
)

// Фразы для помощи пользователем с управлением
const (
	Hello      string = "Приветствую вас на игре «Угадай мелодию»."
	ToStart    string = "Чтобы начать, скажите «Играем»."
	ToStop     string = "Чтобы закончить игру, скажите «Стоп»."
	ToContinue string = "Чтобы продолжить игру, скажите «Продолжить»."
	GoodBye    string = "Пока. Отлично поиграли!"
)

// Для угадываний
const (
	DontGuess    string = "Эх, вы не ^уга`дали^."
	YouGuessText string = "Вы молодец! Угадали!"
	YouGuessTTS  string = "Вы молодец! ^Уга`дали^! <speaker audio=marusia-sounds/game-win-1>"
)

// Для ответов по песням
const (
	IWillSayTheAnswer = "Сейчас скажу ответ. "
	ThatIs            = "Это же: "
)
