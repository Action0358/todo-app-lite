package models

// Todoアイテムを表す構造体
type Todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	Deleted     bool   `json:"deleted"`
}
