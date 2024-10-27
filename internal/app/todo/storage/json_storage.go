package storage

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/AyoubTahir/CLI_Task_Manager/internal/app/todo"
	"github.com/AyoubTahir/CLI_Task_Manager/internal/app/todo/repository"
)

type JSONStorage struct {
	filename string
	mutex    sync.RWMutex
}

func NewJSONStorage(filename string) (repository.TaskRepository, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			return nil, err
		}
	}
	return &JSONStorage{filename: filename}, nil
}

func (s *JSONStorage) readTasks() ([]todo.Task, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	file, err := os.ReadFile(s.filename)
	if err != nil {
		return nil, err
	}

	if len(file) == 0 {
		return []todo.Task{}, nil
	}

	var tasks []todo.Task
	if err := json.Unmarshal(file, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *JSONStorage) writeTasks(tasks []todo.Task) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filename, data, 0644)
}

func (s *JSONStorage) GetAll() ([]todo.Task, error) {
	return s.readTasks()
}

func (s *JSONStorage) Save(task todo.Task) error {
	tasks, err := s.readTasks()
	if err != nil {
		return err
	}
	tasks = append(tasks, task)
	return s.writeTasks(tasks)
}

func (s *JSONStorage) Update(task todo.Task) error {
	tasks, err := s.readTasks()
	if err != nil {
		return err
	}

	for i, t := range tasks {
		if t.ID == task.ID {
			tasks[i] = task
			return s.writeTasks(tasks)
		}
	}
	return nil
}

func (s *JSONStorage) Delete(id int) error {
	tasks, err := s.readTasks()
	if err != nil {
		return err
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return s.writeTasks(tasks)
		}
	}
	return nil
}

func (s *JSONStorage) GetNextID() (int, error) {
	tasks, err := s.readTasks()
	if err != nil {
		return 0, err
	}

	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1, nil
}
