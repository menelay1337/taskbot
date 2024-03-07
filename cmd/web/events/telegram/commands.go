package telegram

import (
	"strings"
	"log"
	"strconv"

	"taskbot/cmd/web/storage"
)

const (
	startCmd = "/start"
	helpCmd = "/help"
	addCmd = "/addtask"
	removeCmd = "/remove"
	tasksCmd = "/tasks"
	pastCmd = "/past"
)

type input struct {
	
}

func (p *Processor) doCmd(text string, chatID int, username string) {
	text = strings.TrimSpace(text)
	
	log.Printf("got new command %s from %s", text, username)

	if isAddCmd(text) {
		return p.saveTask(chatID, input.content, input.days, username)
	}

	switch text:
	case startCmd:
		return p.tg.SendMessage(chatID, msgHello)
	case helpCmd:
		return p.tg.SendMessage(chatID, msgHelp)
	case tasksCmd:
		return p.showTasks(chatID, username)
	case pastCmd:
		return p.pastTasks(chatID, username)
	default:  
		return p.tg.SendMessage(chatID, msgUnknownCommand)
}

func (p *Processor) showTasks(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send tasks", err) }()

	tasks, err := p.storage.Tasks()
	if err != nil && !errors.Is(err, storage.ErrNoSavedTasks {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedTasks) {
		return p.tg.SendMessage(chatID, msgNoSavedTasks)
	}

	return nil
}

//func (p *Processor) pastTasks(chatID int, username string) (err error) {
//	defer func() { err = e.WrapIfErr("can't do command: can't send tasks", err) }()
//
//	tasks, err := p.storage.PastTasks()
//	if err != nil && !errors.Is(err, storage.ErrNoPastTasks {
//		return err
//	}
//
//	if errors.Is(err, storage.ErrNoSavedTasks) {
//		return p.tg.SendMessage(chatID, msgNoPastTasks)
//	}
//
//	return nil
//}

func (p *Processor) saveTask(chatID int, content string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save task", err) }()

	task := &storage.Task {
		Content: content,
	}

	isExists, err := p.storage.IsExists(content)
	if err != nil {
		return err
	}

	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(task); err != nil {
		return err
	}
	
	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}
	return nil
}

func isAddCmd(text string) bool {
	stringSlice := strings.Split(text, " ")
	var answer bool
	days := stringSlice[2]
		
	if (len(stringSlice) != 3 ) {
		return false
	}

	if _, err := strconv.Atoi(days); err != nil  {
		return false
	}
	return true
}
