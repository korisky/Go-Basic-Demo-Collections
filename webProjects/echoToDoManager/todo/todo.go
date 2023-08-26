package todo

import (
	"github.com/labstack/echo/v4"
	"strconv"
	"sync"
	"time"
)

// Todo represents the todo bullet points
type Todo struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	IsComplete bool   `json:"is_complete"`
}

type TodoManager struct {
	todos []Todo
	m     sync.Mutex // operation to be atomic
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}

func NewTodoManager() TodoManager {
	return TodoManager{
		todos: make([]Todo, 0),
		m:     sync.Mutex{},
	}
}

func (tm *TodoManager) GetAll() []Todo {
	return tm.todos
}

func (tm *TodoManager) Create(request CreateTodoRequest) Todo {

	// require lock
	tm.m.Lock()
	defer tm.m.Unlock()

	// create new one
	newTodo := Todo{
		ID:         strconv.FormatInt(time.Now().UnixMilli(), 10),
		Title:      request.Title,
		IsComplete: false,
	}

	// append to collection
	tm.todos = append(tm.todos, newTodo)
	return newTodo
}

func (tm *TodoManager) Complete(ID string) error {
	tm.m.Lock()
	defer tm.m.Unlock()

	var todo *Todo
	var idx = -1

	for i, v := range tm.todos {
		if v.ID == ID {
			idx = i
			todo = &v
		}
	}

	if todo == nil {
		return echo.ErrNotFound
	}

	if todo.IsComplete {
		err := echo.ErrBadRequest
		err.Message = "todo is already completed"
		return err
	}

	// update todo
	tm.todos[idx].IsComplete = true
	return nil
}
