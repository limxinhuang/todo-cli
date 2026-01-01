package todo

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"

func Add(title string) {
	todos, _ := loadTodos()

	newTodo := Todo{
		ID:        GetMaxId(),
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	todos = append(todos, newTodo)
	saveTodos(todos)
	fmt.Println("已添加任务：", title)
}

func List() {
	todos, _ := loadTodos()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\t任务\t状态\t创建时间\t完成时间")
	fmt.Fprintln(w, "")

	for _, t := range todos {
		if t.Completed {
			fmt.Fprintf(w, "%d\t%s\t\033[32m%s\033[0m\t%s\t%s\n", t.ID, t.Title, "V", t.CreatedAt.Format(timeFormat), t.CompletedAt.Format(timeFormat))
		} else {
			fmt.Fprintf(w, "%d\t%s\t\033[31m%s\033[0m\t%s\n", t.ID, t.Title, "X", t.CreatedAt.Format(timeFormat))
		}
	}
	w.Flush()
}

func Completed(id int) {
	todos, _ := loadTodos()
	found := false
	var todo Todo

	for i, t := range todos {
		if t.ID == id {
			todos[i].Completed = true
			todos[i].CompletedAt = time.Now()
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

func DeleteTodo(id int) {
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

func UpdateTitle(id int, title string) {
	todos, _ := loadTodos()
	var oldTitle string

	for i, t := range todos {
		if t.ID == id {
			oldTitle = todos[i].Title
			todos[i].Title = title
			break
		}
	}

	if oldTitle == "" {
		fmt.Printf("任务ID %d 不存在\n", id)
		return
	}

	saveTodos(todos)
	fmt.Printf("任务ID %d 的标题修改成功 %s -> %s\n", id, oldTitle, title)
}
