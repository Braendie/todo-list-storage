package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Braendie/todo-list-storage/internal/models"
)

type Storage struct {
	db *sql.DB
}

func New(storageCon string) *Storage {
	db, err := sql.Open("postgres", storageCon)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping the database:", err)
	}

	return &Storage{db: db}
}

func (s *Storage) CreateTask(ctx context.Context, title string, description string) (int64, error) {
	const op = "storage.pgsql.CreateTask"

	stmt, err := s.db.Prepare(`INSERT INTO tasks (title, description, done) VALUES (?, ?, false)`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, title, description)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetTasks(ctx context.Context) ([]models.Task, error) {
	const op = "storage.pgsql.GetTasks"

	stmt, err := s.db.Prepare(`SELECT id, title, description, done FROM tasks`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	
	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Done); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *Storage) UpdateTask(ctx context.Context, id int64) error {
	const op = "storage.pgsql.UpdateTask"

	stmt, err := s.db.Prepare(`UPDATE tasks SET done = true WHERE id = ?`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteTask(ctx context.Context, id int64) error {
	const op = "storage.pgsql.DeleteTask"

	stmt, err := s.db.Prepare(`DELETE FROM tasks WHERE id = ?`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
