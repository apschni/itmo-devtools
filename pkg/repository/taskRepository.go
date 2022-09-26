package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"homeworkdeliverysystem/dto"
	apperrors "homeworkdeliverysystem/errors"
	"homeworkdeliverysystem/model"
	"log"
	"time"
)

type TaskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (t *TaskRepository) Create(ctx context.Context, task model.Task) (string, error) {
	var id uuid.UUID
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return "", apperrors.NewInternal()
	}

	query := "INSERT INTO tasks " +
		"(id, label, subject, text, deadline, points, closed, teacher_id, file_name, student_id, created_at, updated_at, is_key_point)" +
		" VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id"

	err = t.db.GetContext(ctx, &id, query, newUUID, task.Label, task.Subject, task.Text, task.Deadline, task.Points, task.Closed,
		task.TeacherId, task.FileName, task.StudentId, task.CreatedAt, task.UpdatedAt, task.IsKeyPoint)

	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			log.Printf("Could not create a task with label: %v. Reason: %v\n", task.Label, err.Code.Name())
			return "", apperrors.NewConflict("taskLabel", task.Label)
		}

		log.Printf("Could not create a task with label: %v. Reason: %v\n", task.Label, err)
		return "", apperrors.NewInternal()
	}

	return id.String(), nil
}

func (t *TaskRepository) GetByUserId(ctx context.Context, id uuid.UUID) ([]dto.GetTaskResp, error) {
	var tasks []dto.GetTaskResp

	query := "SELECT t.id, t.label, t.subject, u.full_name AS teacher, t.is_key_point AS keypoint," +
		" t.points, t.closed AS completed, t.deadline " +
		"FROM tasks t JOIN users u on u.id = t.teacher_id WHERE student_id=$1 ORDER BY deadline"

	err := t.db.SelectContext(ctx, &tasks, query, id)
	if err != nil {
		log.Printf("Could not select a task with student_id: %v. Reason: %v\n", id, err)
		return nil, apperrors.NewInternal()
	}

	return tasks, nil
}

func (t *TaskRepository) UpdateFileNameOnMultipleTasks(ctx context.Context, ids pq.StringArray, fileName string) error {
	query := "UPDATE tasks SET file_name=$1, updated_at=$2 WHERE id=$3"

	for _, id := range ids {
		_, err := t.db.ExecContext(ctx, query, fileName, time.Now(), id)
		if err != nil {
			log.Printf("Could not update tasks with ids: %v. Reason: %v\n", ids, err)
			return err
		}
	}
	return nil
}

func (t *TaskRepository) GetFileNameById(ctx context.Context, id string) (string, error) {
	var fileName string

	query := "SELECT file_name FROM tasks WHERE id=$1 LIMIT 1"

	err := t.db.GetContext(ctx, &fileName, query, id)
	if err != nil {
		log.Printf("Could not select a file name from task with id: %v. Reason: %v\n", id, err)
		return "", apperrors.NewInternal()
	}

	return fileName, nil
}

func (t *TaskRepository) Open(ctx context.Context, id uuid.UUID) error {
	query := "UPDATE tasks SET closed=false, updated_at=$1 WHERE id=$2"

	_, err := t.db.ExecContext(ctx, query, time.Now(), id.String())
	if err != nil {
		log.Printf("Could not open id: %v. Reason: %v\n", id, err)
		return err
	}

	return nil
}

func (t *TaskRepository) Close(ctx context.Context, id uuid.UUID) error {
	query := "UPDATE tasks SET closed=true, updated_at=$1 WHERE id=$2"

	_, err := t.db.ExecContext(ctx, query, time.Now(), id.String())
	if err != nil {
		log.Printf("Could not close id: %v. Reason: %v\n", id, err)
		return err
	}

	return nil
}
