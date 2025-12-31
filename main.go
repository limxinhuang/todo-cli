package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"text/tabwriter"
	"time"
)

// Model 层
type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

// 数据文件
const dbFile = "todos.json"

func getDbFile() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, dbFile)
}

// Repository 层

// 读取数据
func loadTodos() ([]Todo, error) {
	if _, err := os.Stat(getDbFile()); os.IsNotExist(err) {
		return []Todo{}, nil
	}

	data, err := os.ReadFile(getDbFile())
	if err != nil {
		return nil, err
	}

	var todos []Todo
	err = json.Unmarshal(data, &todos)
	return todos, err
}

// 保存数据
func saveTodos(todos []Todo) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	// 0644 文件权限 rw-r--r--
	return os.WriteFile(getDbFile(), data, 0644)
}

// Service 层
func add(title string) {
	todos, _ := loadTodos()

	newTodo := Todo{
		ID:        len(todos) + 1,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	todos = append(todos, newTodo)
	saveTodos(todos)
	fmt.Println("已添加任务：", title)
}

func list() {
	todos, _ := loadTodos()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\t任务\t状态\t创建时间")
	fmt.Fprintln(w, "")

	for _, t := range todos {
		status := "X"
		if t.Completed {
			status = "V"
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", t.ID, t.Title, status, t.CreatedAt)
	}
	w.Flush()
}

func completed(id int) {
	todos, _ := loadTodos()
	found := false
	var todo Todo

	for i, t := range todos {
		if t.ID == id {
			todos[i].Completed = true
			found = true
			todo = todos[i]
			break
		}
	}

	if found {
		saveTodos(todos)
		fmt.Printf("任务 %s 已完成\n", todo.Title)
	} else {
		fmt.Printf("任务ID %d 不存在\n", id)
	}
}

func deleteTodo(id int) {
	todos, _ := loadTodos()
	var newTodos []Todo

	for _, t := range todos {
		if t.ID != id {
			newTodos = append(newTodos, t)
		}
	}

	saveTodos(newTodos)
	fmt.Println("任务已删除")
}

// Controller 层

func main() {
	if len(os.Args) < 2 {
		pringHelp()
		return
	}

	cmd := os.Args[1]

	switch cmd {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("错误：请输入任务名称")
			return
		}
		title := os.Args[2]
		add(title)
	case "list":
		list()
	case "done":
		if len(os.Args) < 3 {
			fmt.Println("错误：请输入任务 ID")
			return
		}
		id, _ := strconv.Atoi(os.Args[2])
		completed(id)
	case "del":
		if len(os.Args) < 3 {
			fmt.Println("错误：请输入任务 ID")
			return
		}
		id, _ := strconv.Atoi(os.Args[2])
		deleteTodo(id)
	default:
		pringHelp()
	}
}

func pringHelp() {
	fmt.Println("todo 使用说明:")
	fmt.Println("  add <任务名>   - 添加任务")
	fmt.Println("  list          - 列出所有任务")
	fmt.Println("  done <ID>     - 完成任务")
	fmt.Println("  del <ID>      - 删除任务")
}
