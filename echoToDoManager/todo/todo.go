package todo

import "sync"

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

func NewTodoManager() TodoManager {
	return TodoManager{
		todos: make([]Todo, 0),
		m:     sync.Mutex{},
	}
}

func (tm *TodoManager) GetAll() []Todo {
	return tm.todos
}
