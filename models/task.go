package models

import (
	"time"

	"github.com/go-pg/pg/v10"
)

var TaskModel = Model{
	Name:        "Task",
	Description: "Task Model",
	StructType:  &Task{},
}

type TaskStatus int

// This constants represent the current status of
// the task session.
const (
	PROCESSING TaskStatus = iota // The task is being processed
	SUCCESS                      // The task finished and it was successful
	FAIL                         // The task finished, but it failed.
)

var statusToString = map[TaskStatus]string{
	PROCESSING: "Processing",
	SUCCESS:    "Success",
	FAIL:       "Failure",
}

// The struct represents an task session
type Task struct {
	ID        int        // The unique Number of the task session
	StartDate time.Time  // Task Session start date
	EndDate   time.Time  // Task Session end date
	Status    TaskStatus `pg:",use_zero"` // Task session status
}

func (task *Task) GetStatus() string {
	return statusToString[task.Status]
}

// Creates a new task session.
func NewTaskSession(db *pg.DB, save bool) (*Task, error) {
	newImport := &Task{
		StartDate: time.Now(),
		Status:    PROCESSING,
	}
	if save {
		_, err := db.Model(newImport).Insert()
		if err != nil {
			return nil, err
		}
	}
	return newImport, nil
}

// Returns the latest global task ID
func LatestTaskID(db *pg.DB) (id int, err error) {
	err = db.Model(&Task{}).
		ColumnExpr("max(id)").
		Select(&id)
	return
}

// Returns the latest task ID for a model
func LatestModelTaskID(db *pg.DB, model interface{}) (id int, err error) {
	err = db.Model(model).
		ColumnExpr("max(task_id)").
		Select(&id)
	return
}

// Marks the task session as failed.
// Remember to save this status.
func (task *Task) Failed() {
	task.EndDate = time.Now()
	task.Status = FAIL
}

// Marks the task session as succeeded.
// Remember to save this status.
func (task *Task) Succeeded() {
	task.EndDate = time.Now()
	task.Status = SUCCESS
}

// Save Saves the status of the task session.
func (task *Task) Save(db *pg.DB) (err error) {
	if task.ID == 0 {
		_, err = db.Model(task).Insert()
	} else {
		_, err = db.Model(task).WherePK().Update()
	}
	return
}
