package repository

import "github.com/AyoubTahir/CLI_Task_Manager/internal/app/todo"

type TaskRepository interface {
	Save(task todo.Task) error
	GetAll() ([]todo.Task, error)
	Update(task todo.Task) error
	Delete(id int) error
	GetNextID() (int, error)
}
