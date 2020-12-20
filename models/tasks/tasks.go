package tasks

import (
	"gopkg.in/guregu/null.v4"
	"log"
	"time"
)

const (
	RunScript      string = "RunScript"
	SendParameters        = "SendParameters"
	Reboot                = "Reboot"
	UploadFirmware        = "UploadFirmware"
)

type Task struct {
	Id        int64     `json:"id" db:"id"`
	CpeUuid   string    `json:"cpe_uuid" db:"cpe_uuid"`
	Event     string    `json:"event" db:"event"`
	NotBefore time.Time `json:"not_before" db:"not_before"`
	Task      string    `json:"task" db:"task"`
	Script    string    `json:"script" db:"script"`
	Infinite  bool      `json:"infinite" db:"infinite"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	DoneAt    null.Time `json:"done_at" db:"done_at"`
}

func FilterTasksByEvent(event string, tasksList []Task) []Task {
	var filteredTasks []Task
	for _, task := range tasksList {
		log.Println(task.Event, event)
		log.Println(task.Event == event)
		if task.Event == event {
			filteredTasks = append(filteredTasks, task)
		}
	}

	return filteredTasks
}
