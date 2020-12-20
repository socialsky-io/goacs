package mysql

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"goacs/models/tasks"
	"log"
	"time"
)

type TasksRepository struct {
	db *sqlx.DB
}

func NewTasksRepository(connection *sqlx.DB) TasksRepository {
	return TasksRepository{
		db: connection,
	}
}

func (t *TasksRepository) AddTaskForCPE(task tasks.Task) {
	dialect := goqu.Dialect("mysql")

	query, args, _ := dialect.Insert("tasks").Prepared(true).
		Cols("cpe_uuid", "event", "task", "not_before", "script", "infinite", "created_at").
		Vals(goqu.Vals{
			task.CpeUuid,
			task.Event,
			task.Task,
			task.NotBefore,
			task.Script,
			task.Infinite,
			task.CreatedAt,
		}).ToSQL()

	_, err := t.db.Exec(query, args...)

	if err != nil {
		log.Println("error AddTaskForCPE ", err.Error())
	}

}

func (t *TasksRepository) GetTasksForCPE(cpe_uuid string) []tasks.Task {
	var cpeTasks []tasks.Task
	err := t.db.Select(&cpeTasks, "SELECT * FROM tasks WHERE cpe_uuid=? AND (done_at is null or infinite = true)", cpe_uuid)

	if err != nil {
		log.Println(err.Error())
	}

	return cpeTasks
}

func (t *TasksRepository) GetTasksForCPEWithoutDateCheck(cpe_uuid string) []tasks.Task {
	var cpeTasks []tasks.Task
	_ = t.db.Select(&cpeTasks, "SELECT * FROM tasks WHERE cpe_uuid=? AND done_at is null", cpe_uuid)

	return cpeTasks
}

func (t *TasksRepository) GetAllTasksForCPE(cpe_uuid string) []tasks.Task {
	var cpeTasks []tasks.Task
	_ = t.db.Select(&cpeTasks, "SELECT * FROM tasks WHERE cpe_uuid=?", cpe_uuid)

	return cpeTasks
}

func (t *TasksRepository) DoneTask(task_id int64) {
	_, _ = t.db.Exec("UPDATE tasks SET done_at = ? WHERE id = ?", time.Now(), task_id)

}
