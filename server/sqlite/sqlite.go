package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/Action0358/todo-app-lite/server/models"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStorage struct {
	DB *sql.DB
}

// SQLiteストレージを初期化
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	// データベースファイルのパスを指定
	if dbPath == "" {
		dbPath = "/app/database/todos.db"
	}
	// SQLiteデータベースに接続
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		// 接続に失敗した場合はnilとエラーメッセージを返す
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// テーブルが存在しない場合は作成
	query := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT 0,
		deleted BOOLEAN NOT NULL DEFAULT 0
	);`
	// クエリを実行してテーブルを作成
	_, err = db.Exec(query)
	if err != nil {
		// クエリ実行に失敗した場合はnilとエラーメッセージを返す
		return nil, fmt.Errorf("failed to create table: %v", err)
	}
	// 初期化したSQLiteStorageを返す
	return &SQLiteStorage{DB: db}, nil
}

// データベースのすべてのTodoを取得
func (s *SQLiteStorage) GetAll() ([]models.Todo, error) {
	// todosテーブルから削除されていないアイテムのみを全件取得するクエリを実行
	query := "SELECT id, title, description, completed, deleted FROM todos WHERE deleted = 0"
	// クエリを実行して全てのTodoを取得
	rows, err := s.DB.Query(query)
	if err != nil {
		// クエリ実行時にエラーが発生した場合、nilとエラーメッセージを返す
		return nil, fmt.Errorf("failed to fetch todos: %v", err)
	}
	// 処理が終了またはエラー時にリソース解放を行う
	defer rows.Close()

	// 取得した全てのデータを格納するためのスライスを準備
	var todos []models.Todo
	// テーブル内の次のレコードが存在する間繰り返し実行される
	for rows.Next() {
		// 現在のレコードを格納するための一時的な変数を定義
		var todo models.Todo
		// 現在の行のデータを構造体フィールドにマッピング
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.Deleted)
		if err != nil {
			// データ取得中にエラーが発生した場合、nilとエラーメッセージを返す
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		// スライスtodosに現在のレコードのデータを末尾に追加
		todos = append(todos, todo)
	}
	// 正常の場合、すべてのレコードを格納したスライスを呼び出し元に返す
	return todos, nil
}

// 注: 現在のフロントエンド実装では使用されていない機能です。将来の拡張のために保持。
/*
// ID を指定して Todo を取得
func (s *SQLiteStorage) GetByID(id int) (models.Todo, error) {
	// 指定した ID のレコードを取得するクエリ
	query := "SELECT id, title, description, completed, deleted FROM todos WHERE id = ?"
	// クエリを実行して指定した Todo を取得
	row := s.DB.QueryRow(query, id)

	var todo models.Todo
	// クエリ結果を Todo 構造体にスキャン
	err := row.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.Deleted)
	if err != nil {
		if err == sql.ErrNoRows {
			// 指定した ID が見つからない場合
			return models.Todo{}, fmt.Errorf("failed to retrieve todo: %v", err)
		}
		// その他のエラー
		return models.Todo{}, fmt.Errorf("failed to retrieve todo: %v", err)
	}
	// 正常な場合、todo と nil を返す
	return todo, nil
}
*/

// データベースに新たなTodoを追加
func (s *SQLiteStorage) Create(todo models.Todo) (models.Todo, error) {
	// todosテーブルに title, description, completedを挿入するクエリを実行
	query := "INSERT INTO todos (title, description, completed, deleted) VALUES (?, ?, ?, ?)"
	// クエリを実行してTodoを追加(削除フラグは常に0(false)で初期化)
	result, err := s.DB.Exec(query, todo.Title, todo.Description, todo.Completed, false)
	if err != nil {
		// データ挿入中にエラーが発生した場合、0とエラーメッセージを返す
		return models.Todo{}, fmt.Errorf("failed to insert todo: %v", err)
	}

	// LastInsertId()で挿入されたレコードのIDを取得
	id, err := result.LastInsertId()
	if err != nil {
		// データ取得中にエラーが発生した場合、0とエラーメッセージを返す
		return models.Todo{}, fmt.Errorf("failed to get last insert id: %v", err)
	}
	// 正常の場合、todoとnilを呼び出し元に返す
	todo.ID = int(id)
	return todo, nil
}

// 指定されたIDのTodoを更新する
func (s *SQLiteStorage) Update(id int, todo models.Todo) error {
	// todosテーブルのtitle, description, completedカラムを指定されたidに基づいて更新
	query := "UPDATE todos SET title = ?, description = ?, completed = ?, deleted = ? WHERE id = ?"
	// クエリを実行して更新処理を行う
	result, err := s.DB.Exec(query, todo.Title, todo.Description, todo.Completed, todo.Deleted, id)
	if err != nil {
		// クエリの実行に失敗した場合、エラーを返す
		return fmt.Errorf("failed to update todo: %v", err)
	}
	// 更新された行数を取得
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// 更新行数の取得に失敗した場合、エラーを返す
		return fmt.Errorf("failed to get affected rows: %v", err)
	}
	// 更新された行数が0の場合（指定されたIDが存在しない場合）
	if rowsAffected == 0 {
		// "指定されたIDが見つからない" というエラーを返す
		return sql.ErrNoRows
	}
	// 更新が成功した場合、nilを返す
	return nil
}

// 指定されたIDのTodoを削除する
func (s *SQLiteStorage) Delete(id int) error {
	// 論理削除
	query := "UPDATE todos SET deleted = 1 WHERE id = ?"
	// クエリを実行して削除処理を行う
	result, err := s.DB.Exec(query, id)
	if err != nil {
		// クエリの実行に失敗した場合、エラーを返す
		return fmt.Errorf("failed to delete todo: %v", err)
	}
	// 削除された行数を取得
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// 削除行数の取得に失敗した場合、エラーを返す
		return fmt.Errorf("failed to get affected rows: %v", err)
	}
	// 削除された行数が0の場合（指定されたIDが存在しない場合）
	if rowsAffected == 0 {
		// "指定されたIDが見つからない"というエラーを返す
		return sql.ErrNoRows
	}
	// 削除が成功した場合、nilを返す
	return nil
}
