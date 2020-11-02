package mysql

import (
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
