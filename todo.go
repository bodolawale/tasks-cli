package main

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"
)

const fileName = "todos.csv"

type Todo struct {
	ID          int
	Description string
	Completed   bool
	CreatedAt   time.Time
}

func saveTodo(todo Todo) error {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		strconv.Itoa(todo.ID),
		todo.Description,
		strconv.FormatBool(todo.Completed),
		todo.CreatedAt.Format(time.RFC3339),
	}

	return writer.Write(record)
}

func loadTodos() ([]Todo, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var todos []Todo

	for _, row := range rows {
		id, _ := strconv.Atoi(row[0])
		completed, _ := strconv.ParseBool(row[2])
		createdAt, _ := time.Parse(time.RFC3339, row[3])

		todo := Todo{
			ID:          id,
			Description: row[1],
			Completed:   completed,
			CreatedAt:   createdAt,
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func getNextID(todos []Todo) int {
	maxID := 0

	for _, todo := range todos {
		if todo.ID > maxID {
			maxID = todo.ID
		}
	}

	return maxID + 1
}

func saveAllTodos(todos []Todo) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, todo := range todos {
		record := []string{
			strconv.Itoa(todo.ID),
			todo.Description,
			strconv.FormatBool(todo.Completed),
			todo.CreatedAt.Format(time.RFC3339),
		}

		err := writer.Write(record)
		if err != nil {
			return err
		}
	}
	return nil
}
