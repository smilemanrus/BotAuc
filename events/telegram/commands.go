package telegram

import (
	"BotAuc/lib/e"
	"context"
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
		return p.sendAucData(chatID)
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

func (p *Processor) sendAucData(chatID int) error {
	msg := ""
	err := p.storage.GetFutureAucs(context.Background(), &msg)
	if err != nil {
		err = e.Wrap("can't get auc info", err)
		return err
	}
	return p.tg.SendMessage(chatID, msg)
}
