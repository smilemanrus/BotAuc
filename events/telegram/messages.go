package telegram

const msgHelp = `Этот бот будет оповещать вас о предстоящих аукционах в фед офисе.
Уведомления будут приходить в момент публикации аукциона, его обновления, за 2 часа и час до начала торгов.
Отправьте "/auc" чтобы самостоятельно проверить наличие предстоящих аукционов`

const msgStart = `Привет! 
` + msgHelp

const (
	msgUnknownCmd  = "Неизвестная команда"
	msgNoAuc       = "Предстоящих аукционов нет"
	msgNoFunc      = "Эта функциональность ещё не реализована"
	msgSubscrAuc   = "Оповещения об аукционах включены"
	msgUnSubscrAuc = "Оповещения об аукционах выключены"
)
