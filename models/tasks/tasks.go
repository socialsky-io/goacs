package tasks

import (
	"database/sql"
	"time"
)

const (
	RunScript string = "RunScript"
	Provision        = "Provision"
	Reboot           = "Reboot"
)

type Task struct {
	Id        int64        `json:"id" db:"id"`
	CpeUuid   string       `json:"cpe_uuid" db:"cpe_uuid"`
	Event     string       `json:"event" db:"event"`
	NotBefore time.Time    `json:"not_before" db:"not_before"`
	Task      string       `json:"task" db:"task"`
	Script    string       `json:"script" db:"script"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	DoneAt    sql.NullTime `json:"done_at" db:"done_at"`
}

func FilterTasksByEvent(event string, tasksList []Task) []Task {
	var filteredTasks []Task
	for _, task := range tasksList {
		if task.Event == event {
			filteredTasks = append(filteredTasks, task)
		}
	}

	return filteredTasks
}
