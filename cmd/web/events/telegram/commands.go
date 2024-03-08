package telegram

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	//"net/url"

	"taskbot/cmd/web/storage"
	"taskbot/internal/e"
)

var isLogin bool

const (
	startCmd    = "/start"
	addCmd      = "/addtask"
	removeCmd   = "/remove"
	tasksCmd    = "/tasks"
	completeCmd = "/complete"
	authCmd     = "/auth"
	registerCmd = "/register"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command %s from %s", text, username)

	if isAddCmd(text) {
		return p.saveTask(chatID, text, username)
	}

	switch text {
	case registerCmd:
		p.Register(chatID, username)
	case startCmd:
		p.Auth(chatID, username)
		return p.tg.SendMessage(chatID, msgHello)
	case tasksCmd:
		return p.showTasks(chatID, username)
	case completeCmd:
		return p.completeTask(chatID, text, username)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
	return nil
}
func (p *Processor) Register(chatid int, username string) (err error) {
	pretendent := &storage.User{
		Username: username,
		Chatid:   chatid,
	}
	_, err = p.storage.RetrieveUser(pretendent)
	if err != nil {
		errInSaving := p.storage.SaveUser(pretendent)
		if errInSaving != nil {
			return p.tg.SendMessage(chatid, msgPlsRegister)
		}
		isLogin = true
	}
	return p.tg.SendMessage(chatid, msgUserExist)
}
func (p *Processor) Auth(chatid int, username string) (err error) {
	//check
	pretendent := &storage.User{
		Username: username,
		Chatid:   chatid,
	}
	_, err = p.storage.RetrieveUser(pretendent)
	if err != nil {
		return p.tg.SendMessage(chatid, msgPlsRegister)
	}
	return p.tg.SendMessage(chatid, msgUserExist)
}
func (p *Processor) showTasks(chatID int, username string) (err error) {
	if isLogin {
		defer func() { err = e.WrapIfErr("can't do command: can't send tasks", err) }()

		tasks, err := p.storage.Tasks()
		if err != nil && !errors.Is(err, storage.ErrNoSavedTasks) {
			return err
		}

		if errors.Is(err, storage.ErrNoSavedTasks) {
			return p.tg.SendMessage(chatID, msgNoSavedTasks)
		}

		taskListText := "Task List:\n"
		for _, task := range tasks {
			completedStatus := "Not Completed"
			if task.Completed {
				completedStatus = "Completed"
			}
			taskListText += fmt.Sprintf("- Task %d: %s (Created: %s, %s)\n", task.ID, task.Content, task.Created.Format("2006-01-02 15:04:05"), completedStatus)
		}

		p.tg.SendMessage(chatID, taskListText)

		return nil
	}
	return p.tg.SendMessage(chatID, msgPlsRegister)

}

func (p *Processor) completeTask(chatID int, text string, username string) error {
	if isLogin {
		textSlice := strings.Split(text, "")
		textSlice[1] = strings.Trim(textSlice[1], "\"")
		_, err := strconv.Atoi(textSlice[1])
		if len(textSlice) == 2 && err == nil {
			id, _ := strconv.Atoi(textSlice[1])
			p.storage.Complete(id)
			p.tg.SendMessage(chatID, msgCompleted)
		} else {
			return p.tg.SendMessage(chatID, msgIncorrectInput)
		}

		return nil
	}
	return p.tg.SendMessage(chatID, msgPlsRegister)

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

func (p *Processor) saveTask(chatID int, text string, username string) (err error) {
	if isLogin {
		defer func() { err = e.WrapIfErr("can't do command: save task", err) }()
		textSlice := strings.Split(text, "")
		var content string
		if len(textSlice) == 2 {
			content = strings.Trim(textSlice[1], "\"")
		} else {
			return p.tg.SendMessage(chatID, msgIncorrectInput)
		}
		task := &storage.Task{
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
	return p.tg.SendMessage(chatID, msgPlsRegister)

}

func isAddCmd(text string) bool {
	textSlice := strings.Split(text, " ")

	if len(textSlice) == 2 && textSlice[0] == addCmd {
		return true
	}

	return false
}
