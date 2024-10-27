package main

import (
	"log"

	"github.com/AyoubTahir/CLI_Task_Manager/internal/app/todo/storage"
	"github.com/AyoubTahir/CLI_Task_Manager/internal/app/todo/ui"
)

func main() {
	storage, err := storage.NewJSONStorage("tasks.json")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	cli := ui.NewCLI(storage)
	cli.Run()
}
