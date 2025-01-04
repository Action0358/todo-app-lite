package storage

import (
	"database/sql"
	"fmt"

	"github.com/Action0358/todo-app-lite/server/models"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStorage struct {
	DB *sql.DB
}

// SQLite ストレージを初期化
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	// データベースファイルのパスを指定
	if dbPath == "" {
		dbPath = "/app/database/todos.db"
	}
	// SQLite データベースに接続
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		// 接続に失敗した場合は nil とエラーメッセージを返す
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// テーブルが存在しない場合は作成
	query := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		completed BOOLEAN NOT NULL
	);`
	// クエリを実行してテーブルを作成
	_, err = db.Exec(query)
	if err != nil {
		// クエリ実行に失敗した場合は nil とエラーメッセージを返す
		return nil, fmt.Errorf("failed to create table: %v", err)
	}
	// 初期化した SQLiteStorage を返す
	return &SQLiteStorage{DB: db}, nil
}

// データベースのすべての Todo を取得
func (s *SQLiteStorage) GetAll() ([]models.Todo, error) {
	// todos テーブルから id, title, completed を全件取得するクエリを実行
	query := "SELECT id, title, description, completed FROM todos"
	// クエリを実行して全ての Todo を取得
	rows, err := s.DB.Query(query)
	if err != nil {
		// クエリ実行時にエラーが発生した場合、nil とエラーメッセージを返す
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
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed)
		if err != nil {
			// データ取得中にエラーが発生した場合、nil とエラーメッセージを返す
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		// スライス todos に現在のレコードのデータを末尾に追加
		todos = append(todos, todo)
	}
	// 正常の場合、すべてのレコードを格納したスライスを呼び出し元に返す
	return todos, nil
}

// ID を指定して Todo を取得
func (s *SQLiteStorage) GetByID(id int) (models.Todo, error) {
	// 指定した ID のレコードを取得するクエリ
	query := "SELECT id, title, description, completed FROM todos WHERE id = ?"
	// クエリを実行して指定した Todo を取得
	row := s.DB.QueryRow(query, id)

	var todo models.Todo
	// クエリ結果を Todo 構造体にスキャン
	err := row.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed)
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

// データベースに新たな Todo を追加
func (s *SQLiteStorage) Add(todo models.Todo) (int64, error) {
	// todos テーブルに title, description, completed を挿入するクエリを実行
	query := "INSERT INTO todos (title, description, completed) VALUES (?, ?, ?)"
	// クエリを実行して Todo を追加
	result, err := s.DB.Exec(query, todo.Title, todo.Description, todo.Completed)
	if err != nil {
		// データ挿入中にエラーが発生した場合、 0 とエラーメッセージを返す
		return 0, fmt.Errorf("failed to insert todo: %v", err)
	}

	// LastInsertId()で挿入されたレコードの ID を取得
	id, err := result.LastInsertId()
	if err != nil {
		// データ取得中にエラーが発生した場合、 0 とエラーメッセージを返す
		return 0, fmt.Errorf("failed to get last insert id: %v", err)
	}
	// 正常の場合、id と nil を呼び出し元に返す
	return id, nil
}

// 指定された ID の Todo を更新する
func (s *SQLiteStorage) Update(id int, todo models.Todo) error {
	// todos テーブルの title, description, completed カラムを指定された id に基づいて更新
	query := "UPDATE todos SET title = ?, description = ?, completed = ? WHERE id = ?"
	// クエリを実行して更新処理を行う
	result, err := s.DB.Exec(query, todo.Title, todo.Description, todo.Completed, id)
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
	// 更新された行数が 0 の場合（指定された ID が存在しない場合）
	if rowsAffected == 0 {
		// "指定された ID が見つからない" というエラーを返す
		return sql.ErrNoRows
	}
	// 更新が成功した場合、nil を返す
	return nil
}

// 指定された ID の Todo を削除
func (s *SQLiteStorage) Delete(id int) error {
	// 削除対象の id に基づいて、todos テーブルから該当行を削除するクエリ
	query := "DELETE FROM todos WHERE id = ?"
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
	// 削除された行数が 0 の場合（指定された ID が存在しない場合）
	if rowsAffected == 0 {
		// "指定された ID が見つからない" というエラーを返す
		return sql.ErrNoRows
	}
	// 削除が成功した場合、nil を返す
	return nil
}
