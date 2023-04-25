package telegram

import (
	"log"
	"strings"
)

const (
	AucCmd   = "/auc"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCMD(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)
	switch text {
	case AucCmd:
		return p.sendNoFunc(chatID)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendStart(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCmd)
	}

}

func (p *Processor) sendNoFunc(chatID int) error {
	return p.tg.SendMessage(chatID, msgNoFunc)
}
func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendStart(chatID int) error {
	return p.tg.SendMessage(chatID, msgStart)
}
