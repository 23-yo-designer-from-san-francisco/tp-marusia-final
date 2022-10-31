package answer

// Управляющие фразы
const (
	Playem   string = "играем"
	Play     string = "играть"
	Continue string = "продолжить"
	Greeting string = "приветствие"
	Begin    string = "начать"
)

// Ответы: Не знаю
const (
	IDontKnow  string = "я не знаю"
	DontKnow   string = "не знаю"
	No         string = "нет"
	CantGuess         = "не могу угадать"
	ICantGuess        = "я не могу угадать"
)

// Для начала и конца игры
const (
	Hello   string = `Приветствую вас на игре «Угадай мелодию». Я Маруся, буду сегодня вашим ведущим.`
	GoodBye string = `Пока. Отлично поиграли`
)

// Для не тех фраз в определённых местах
const (
	AlreadyPlaying    = "Вы уже играете, забыли?"
	UnknownCommand    = "Я вас не совсем поняла"
	UnknownCommandTTS = "Я вас «НЕ СОВСЕМ» поняла"
)

// Фразы для помощи пользователем с управлением
const (
	ToStart    string = `Чтобы начать, скажите «Играем». `
	ToContinue string = `Чтобы продолжить игру, скажите «Продолжить». `
	ToStop     string = `Чтобы закончить игру, скажите «Стоп». `
)

// Для угадываний
const (
	DontGuess string = "Эх, вы не ^уга`дали^."
	YouGuess  string = "Вы молодец! ^Уга`дали^! <speaker audio=marusia-sounds/game-win-1>"
)

// Для ответов по песням
const (
	IWillSayTheAnswer = "Сейчас скажу ответ."
	ThatIs            = "Это же: "
)
