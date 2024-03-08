package telegram

const msgHelp = `I can save and keep your tasks. In any time you can watch past and present tasks.

In order to save the task, just send the command /add and after that send content and then send the number of days to complete this task.

please enter data with the pattern /command "content"
`

const msgHello = "Hi there \n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command"
	msgNoSavedTasks = "You have no saved task"
	msgNoPastTasks = "There are no past tasks"
	msgSaved = "Task is saved."
	msgRemoved = "Task is removed."
	msgCompleted = "Task is completed."
	msgAlreadyExists = "You have already have this task in your list"

	msgIncorrectInput = "Incorrept input"
)
