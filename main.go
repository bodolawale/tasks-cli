package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

func main() {
	// fmt.Println(len(os.Args))

	if len(os.Args) < 2 {
		fmt.Println("Usage: tasks <command>")
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: tasks add <description>")
			return
		}
		description := strings.Join(os.Args[2:], " ")
		addTask(description)
	case "list":
		showAll := false
		if len(os.Args) == 3 {
			if os.Args[2] == "-a" {
				showAll = true
			} else {
				fmt.Println("Invalid Argument: Usage: list -a")
				return
			}
		}
		listTasks(showAll)
	case "complete":
		if len(os.Args) != 3 {
			fmt.Println("Usage: tasks complete <taskID>")
			return
		}
		taskID, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("%v must be valid task id, an integer\n", os.Args[2])
			return
		}
		completeTask(taskID)
	case "uncomplete":
		if len(os.Args) != 3 {
			fmt.Println("Usage: tasks uncomplete <taskID>")
			return
		}
		taskID, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("%v must be a valid integer\n", os.Args[2])
			return
		}
		markTaskAsIncomplete(taskID)
	case "delete":
		if len(os.Args) != 3 {
			fmt.Printf("Usage: tasks delete <taskID>\n")
			return
		}
		taskID, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("%v must be a valid integer\n", os.Args[2])
			return
		}
		deleteTask(taskID)
	default:
		fmt.Println("Unknown command:", command)
	}

}

func addTask(description string) {
	todos, err := loadTodos()
	if err != nil {
		if os.IsNotExist(err) {
			todos = []Todo{}
		} else {
			fmt.Println("Error loading todos:", err)
			return
		}
	}

	nextID := getNextID(todos)

	todo := Todo{
		ID:          nextID,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
	}

	err = saveTodo(todo)
	if err != nil {
		fmt.Println("Error saving todo:", err)
		return
	}

	fmt.Println("Task added successfully")
}

func listTasks(showAll bool) {
	todos, err := loadTodos()
	if err != nil {
		fmt.Println("Error listing todos:", err)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	if showAll {
		fmt.Fprintln(w, "ID\tTask\tCreated\tDone")
	} else {
		fmt.Fprintln(w, "ID\tTask\tCreated")
	}

	for _, todo := range todos {
		if !showAll && todo.Completed {
			continue
		}

		if showAll {
			doneValue := "pending"
			if todo.Completed {
				doneValue = "done"
			}
			fmt.Fprintf(
				w,
				"%d\t%s\t%s\t%s\n",
				todo.ID,
				todo.Description,
				timeAgo(todo.CreatedAt),
				doneValue,
			)
		} else {
			fmt.Fprintf(
				w,
				"%d\t%s\t%s\n",
				todo.ID,
				todo.Description,
				timeAgo(todo.CreatedAt),
			)
		}

	}

	w.Flush()
}

func completeTask(taskID int) {
	todos, err := loadTodos()
	if err != nil {
		if os.IsNotExist(err) {
			todos = []Todo{}
		} else {
			fmt.Println("Error loading todos:", err)
			return
		}
	}

	exists := false
	for i, todo := range todos {
		if todo.ID == taskID {
			if todo.Completed {
				fmt.Println("Task is already complete")
				return
			}
			todos[i].Completed = true
			exists = true
			break
		}
	}

	if !exists {
		fmt.Printf("Todo with id %d not found\n", taskID)
		return
	}

	err = saveAllTodos(todos)
	if err != nil {
		fmt.Println("Error saving all todos:", err)
		return
	}

	fmt.Println("Todo marked as completed")
}

func markTaskAsIncomplete(taskID int) {
	todos, err := loadTodos()
	if err != nil {
		fmt.Println("Error loading todos:", err)
		return
	}

	exists := false
	for i, todo := range todos {
		if todo.ID == taskID {
			if !todo.Completed {
				fmt.Println("Todo is already incomplete")
				return
			}
			todos[i].Completed = false
			exists = true
			break
		}
	}

	if !exists {
		fmt.Printf("Todo with id %d not found\n", taskID)
		return
	}

	err = saveAllTodos(todos)
	if err != nil {
		fmt.Println("Error saving todos:", err)
		return
	}

	fmt.Println("Todo marked as incomplete")
}

func deleteTask(taskID int) {
	todos, err := loadTodos()
	if err != nil {
		fmt.Println("Error loading todos:", err)
		return
	}

	newTodos := []Todo{}
	found := false
	var todoToDelete Todo

	for _, todo := range todos {
		if todo.ID == taskID {
			found = true
			todoToDelete = todo
			continue
		}
		newTodos = append(newTodos, todo)
	}

	if !found {
		fmt.Printf("Todo with id %d not found\n", taskID)
		return
	}

	err = saveAllTodos(newTodos)
	if err != nil {
		fmt.Println("Error saving todos:", err)
		return
	}

	fmt.Printf("Deleted: %s\n", todoToDelete.Description)
}

func timeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		return "a few seconds ago"
	}

	if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	}

	if duration < time.Hour*24 {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}

	days := int(duration.Hours() / 24)
	if days == 1 {
		return "1 day ago"
	}

	return fmt.Sprintf("%d days ago", days)
}
