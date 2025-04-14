package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Action0358/todo-app-lite/server/models"
	"github.com/Action0358/todo-app-lite/server/sqlite"
)

// SQLiteストレージのポインタ変数を宣言
var handlerStorage *sqlite.SQLiteStorage

// 外部からストレージを設定する関数
func SetStorage(storage *sqlite.SQLiteStorage) {
	handlerStorage = storage
}

// `/todos`エンドポイントを処理
func TodosHandlers(w http.ResponseWriter, r *http.Request) {
	// レスポンスヘッダにJSONコンテンツタイプを設定
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// 全てのTodoを取得するが、todosをレスポンスに返さない
		todos, err := handlerStorage.GetAll()
		// 取得に失敗した場合、HTTP500エラーを返す
		if err != nil {
			http.Error(w, "Failed to retrieve todos", http.StatusInternalServerError)
			return
		}
		// 取得したTodoの一覧をJSON形式でエンコードし、レスポンスとして送信（元のデータ構造 -> JSON形式に変換）
		json.NewEncoder(w).Encode(todos)

	case http.MethodPost:
		// Todoの追加
		var newTodo models.Todo
		// リクエストボディからJSONデータをデコードして、newTodo構造体に格納
		err := json.NewDecoder(r.Body).Decode(&newTodo)
		// デコードエラー時、クライアントに400Bad Requestを返す
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		// 新規作成時は常に削除フラグをfalseに設定
		newTodo.Completed = false
		newTodo.Deleted = false
		// ストレージに新しいTodoを追加
		createdTodo, err := handlerStorage.Create(newTodo)
		// 追跡処理でエラーが発生した場合、500Internal Server Errorを返す
		if err != nil {
			http.Error(w, "Failed to add todo", http.StatusInternalServerError)
			return
		}
		// バリデーションチェック
		if newTodo.Title == "" || newTodo.Description == "" {
			http.Error(w, "Please enter your information", http.StatusBadRequest)
			return
		}
		// ステータスコードを201Createdに設定
		w.WriteHeader(http.StatusCreated)
		// 追加されたTodoをJSON形式で応答（元のデータ構造 -> JSON形式に変換)
		json.NewEncoder(w).Encode(createdTodo)
	}
}

// `/todos/{id}`エンドポイントを処理
func TodoHandlers(w http.ResponseWriter, r *http.Request) {
	// レスポンスヘッダにJSONコンテンツタイプを設定
	w.Header().Set("Content-Type", "application/json")

	// リクエストURLからIDを抽出
	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	// 文字列を整数に変換
	id, err := strconv.Atoi(idStr)
	// 変換に失敗した場合、またはIDが空の場合、エラーレスポンスを返す
	if err != nil || idStr == "" {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	// 注: GETメソッドは現在のフロントエンド実装では使用されていない機能です。将来の拡張のために保持。
	/*
		case http.MethodGet:
			// ストレージ内の指定されたIDのTodoを取得
			todo, err := handlerStorage.GetByID(id)
			// 取得に失敗した場合、HTTP 404 not foundを返す
			if err != nil {
				http.Error(w, "Todo not found", http.StatusNotFound)
			}
			// 取得されたTodoをJSON形式でエンコードし、レスポンスとして送信（元のデータ構造 -> JSON形式に変換）
			json.NewEncoder(w).Encode(todo)
	*/

	case http.MethodPut:
		// 指定されたIDのTodoを更新
		var updatedTodo models.Todo
		// リクエストボディからJSONデータをデコードして、updateTodo構造体に格納（JSON形式 -> 元のデータ構造に変換）
		err := json.NewDecoder(r.Body).Decode(&updatedTodo)
		// デコードエラー時、クライアントに400 Bad Requestを返す
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		// 更新するTodoにURLから取得したIDを設定
		updatedTodo.ID = id
		// ストレージ内の指定されたIDのTodoを新しいデータで更新
		err = handlerStorage.Update(id, updatedTodo)
		// 削除に失敗した場合、HTTP 404 not foundを返す
		if err != nil {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		// バリデーションチェック
		if updatedTodo.Title == "" || updatedTodo.Description == "" {
			http.Error(w, "Please enter your information", http.StatusBadRequest)
			return
		}
		// 更新されたTodoをJSON形式でエンコードし、レスポンスとして送信（元のデータ構造 -> JSON形式に変換）
		json.NewEncoder(w).Encode(updatedTodo)

	case http.MethodDelete:
		// 指定されたIDのTodoを論理削除
		err := handlerStorage.Delete(id)
		// 削除に失敗した場合、HTTP 404 not foundを返す
		if err != nil {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		// HTTP 204 No Contentを返す
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
