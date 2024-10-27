package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/AyoubTahir/CLI_Task_Manager/internal/app/todo"
	"github.com/AyoubTahir/CLI_Task_Manager/internal/app/todo/repository"
)

type CLI struct {
	repository repository.TaskRepository
	reader     *bufio.Reader
}

func NewCLI(repo repository.TaskRepository) *CLI {
	return &CLI{
		repository: repo,
		reader:     bufio.NewReader(os.Stdin),
	}
}

func (c *CLI) readInput(prompt string) string {
	fmt.Print(prompt)
	input, _ := c.reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func (c *CLI) addTask() error {
	title := c.readInput("Enter task title: ")
	description := c.readInput("Enter task description: ")

	id, err := c.repository.GetNextID()
	if err != nil {
		return err
	}

	task := todo.Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      false,
		CreatedAt:   time.Now(),
	}

	return c.repository.Save(task)
}

func (c *CLI) listTasks() error {
	tasks, err := c.repository.GetAll()
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found!")
		return nil
	}

	fmt.Println("\nYour Tasks:")
	fmt.Println("----------------------------------------")
	for _, task := range tasks {
		status := "[ ]"
		if task.Status {
			status = "[✓]"
		}
		fmt.Printf("%d. %s %s\n", task.ID, status, task.Title)
		fmt.Printf("   Description: %s\n", task.Description)
		fmt.Printf("   Created: %s\n", task.CreatedAt.Format("2006-01-02 15:04:05"))
		if task.Status {
			fmt.Printf("   Completed: %s\n", task.CompletedAt.Format("2006-01-02 15:04:05"))
		}
		fmt.Println("----------------------------------------")
	}
	return nil
}

func (c *CLI) completeTask() error {
	idStr := c.readInput("Enter task ID to complete: ")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid ID")
	}

	tasks, err := c.repository.GetAll()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if task.ID == id {
			task.Status = true
			task.CompletedAt = time.Now()
			return c.repository.Update(task)
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

func (c *CLI) deleteTask() error {
	idStr := c.readInput("Enter task ID to delete: ")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid ID")
	}

	return c.repository.Delete(id)
}

func (c *CLI) Run() {
	for {
		fmt.Println("\n=== Task Manager ===")
		fmt.Println("1. Add Task")
		fmt.Println("2. List Tasks")
		fmt.Println("3. Complete Task")
		fmt.Println("4. Delete Task")
		fmt.Println("5. Exit")

		choice := c.readInput("Choose an option: ")

		var err error
		switch choice {
		case "1":
			err = c.addTask()
		case "2":
			err = c.listTasks()
		case "3":
			err = c.completeTask()
		case "4":
			err = c.deleteTask()
		case "5":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid option! Please try again.")
			continue
		}

		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}
