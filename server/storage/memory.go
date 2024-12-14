package storage

import (
	"fmt"

	"github.com/Action0358/todo-app-lite/server/models"
)

type MemoryStorage struct {
	Todos []models.Todo
}

// グローバル変数を宣言
var instance *MemoryStorage

// MemoryStorage のインスタンスを取得または作成する関数
func NewMemoryStorage() *MemoryStorage {
	if instance == nil {
		// 初回呼び出し時にインスタンス作成
		instance = &MemoryStorage{
			Todos: []models.Todo{},
		}
	}
	// 作成された（または既存の）インスタンスを呼び出し元に返す
	return instance
}

// ストレージ内のすべての Todo を取得
func (s *MemoryStorage) GetAll() ([]models.Todo, error) {
	return s.Todos, nil
}

// ID を 指定して Todo を取得するメソッド
func (s *MemoryStorage) GetByID(id int) (models.Todo, error) {
	// Todo スライスをループして、指定された ID の Todo を探す
	for _, todo := range s.Todos {
		if todo.ID == id {
			return todo, nil
		}
	}
	// 指定した ID が見つからなかった場合、エラーを返す
	return models.Todo{}, fmt.Errorf("Todo with ID %d not found", id)
}

// 新しい Todo を追加
func (s *MemoryStorage) Add(todo models.Todo) error { // todo = newTodo のコピー
	s.Todos = append(s.Todos, todo)
	return nil
}

// 指定された ID の Todo を更新するメソッド
func (s *MemoryStorage) Update(id int, todo models.Todo) error { // todo = updateTodo のコピー
	// Todo スライスをループして、指定された ID の Todo を探す
	for i, t := range s.Todos {
		if t.ID == id {
			s.Todos[i] = todo
			return nil
		}
	}
	// 指定した ID が見つからなかった場合、エラーを返す
	return fmt.Errorf("Todo with ID %d not found", id)
}

// 指定した ID の Todo を削除するメソッド
func (s *MemoryStorage) Delete(id int) error {
	for i, t := range s.Todos {
		if t.ID == id {
			s.Todos = append(s.Todos[:i], s.Todos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Todo with ID %d not found", id)
}
