package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Action0358/todo-app-lite/server/models"
	"github.com/Action0358/todo-app-lite/server/sqlite"
)

// SQLite ストレージのポインタ変数を宣言
var handlerStorage *sqlite.SQLiteStorage

// 外部からストレージを設定する関数
func SetStorage(storage *sqlite.SQLiteStorage) {
	handlerStorage = storage
}

// `/todos` エンドポイントを処理
func TodosHandlers(w http.ResponseWriter, r *http.Request) {
	// レスポンスヘッダに JSON コンテンツタイプを設定
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// 全ての Todo を取得するが、todosをレスポンスに返さない
		todos, err := handlerStorage.GetAll()
		// 取得に失敗した場合、 HTTP 500 エラーを返す
		if err != nil {
			http.Error(w, "Failed to retrieve todos", http.StatusInternalServerError)
			return
		}
		// 取得した Todo の一覧を JSON 形式でエンコードし、レスポンスとして送信（元のデータ構造 -> JSON 形式に変換）
		json.NewEncoder(w).Encode(todos)

	case http.MethodPost:
		// Todo の追加
		var newTodo models.Todo
		// リクエストボディから JSON データをデコードして、 newTodo 構造体に格納
		err := json.NewDecoder(r.Body).Decode(&newTodo)
		// デコードエラー時、クライアントに 400 Bad Request を返す
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		// ストレージに新しい Todo を追加するが、レスポンスにidを返さない
		createdTodo, err := handlerStorage.Create(newTodo)
		// 追跡処理でエラーが発生した場合、 500 Internal Server Error を返す
		if err != nil {
			http.Error(w, "Failed to add todo", http.StatusInternalServerError)
			return
		}
		// バリデーションチェック
		if newTodo.Title == "" || newTodo.Description == "" {
			http.Error(w, "Please enter your information", http.StatusBadRequest)
			return
		}
		// ステータスコードを 201 Created に設定
		w.WriteHeader(http.StatusCreated)
		// 追加された Todo を JSON で応答(元のデータ構造 -> JSON 形式に変換)
		json.NewEncoder(w).Encode(createdTodo)
	}
}

// `/todos/{id}` エンドポイントを処理
func TodoHandlers(w http.ResponseWriter, r *http.Request) {
	// レスポンスヘッダに JSON コンテンツタイプを設定
	w.Header().Set("Content-Type", "application/json")

	// リクエスト URL から ID を抽出
	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	// 文字列を整数に変換
	id, err := strconv.Atoi(idStr)
	// 変換に失敗した場合、または ID が空の場合、エラーレスポンスを返す
	if err != nil || idStr == "" {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	// 注: GET メソッドは現在のフロントエンド実装では使用されていない機能です。将来の拡張のために保持。
    /*
	case http.MethodGet:
		// ストレージ内の指定された ID の Todo を取得
		todo, err := handlerStorage.GetByID(id)
		// 取得に失敗した場合、 HTTP 404 not found を返す
		if err != nil {
			http.Error(w, "Todo not found", http.StatusNotFound)
		}
		// 取得された Todo を JSON 形式でエンコードし、レスポンスとして送信（元のデータ構造 -> JSON 形式に変換）
		json.NewEncoder(w).Encode(todo)
	*/

	case http.MethodPut:
		// 指定された ID の Todo を更新
		var updatedTodo models.Todo
		// リクエストボディから JSON データをデコードして、updateTodo 構造体に格納（JSON 形式 -> 元のデータ構造に変換）
		err := json.NewDecoder(r.Body).Decode(&updatedTodo)
		// デコードエラー時、クライアントに400 Bad Requestを返す
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		// 更新する Todo に URL から取得した ID を設定
		updatedTodo.ID = id
		// ストレージ内の指定された ID の Todo を新しいデータで更新
		err = handlerStorage.Update(id, updatedTodo)
		// 削除に失敗した場合、 HTTP 404 not found を返す
		if err != nil {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		// バリデーションチェック
		if updatedTodo.Title == "" || updatedTodo.Description == "" {
			http.Error(w, "Please enter your information", http.StatusBadRequest)
			return
		}
		// 更新された Todo を JSON 形式でエンコードし、レスポンスとして送信（元のデータ構造 -> JSON 形式に変換）
		json.NewEncoder(w).Encode(updatedTodo)

	case http.MethodDelete:
		// 指定された ID の Todo を削除
		err := handlerStorage.Delete(id)
		// 削除に失敗した場合、 HTTP 404 not found を返す
		if err != nil {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
