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
	switch r.Method {
	case http.MethodGet:
		// 全ての Todo を取得するが、todosをレスポンスに返さない
		_, err := handlerStorage.GetAll()
		// 取得に失敗した場合、 HTTP 500 エラーを返す
		if err != nil {
			http.Error(w, "Failed to retrieve todos", http.StatusInternalServerError)
			return
		}
		// レスポンスヘッダに JSON コンテンツタイプを設定
		w.Header().Set("Content-Type", "application/json")
		// レスポンスに成功メッセージを返す
		response := map[string]string{
			"message": "Todos retrieved successfully",
		}
		// 取得した Todo の一覧を JSON 形式でエンコードし、レスポンスとして送信（元のデータ構造 -> JSON 形式に変換）
		json.NewEncoder(w).Encode(response)

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
		_, err = handlerStorage.Add(newTodo)
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
		// レスポンスヘッダに JSON コンテンツタイプを設定
		w.Header().Set("Content-Type", "application/json")
		// ステータスコードを 201 Created に設定
		w.WriteHeader(http.StatusCreated)
		response := map[string]string{
			"message": "Todo added successfully",
		}
		// 追加された Todo を JSON で応答(元のデータ構造 -> JSON 形式に変換)
		json.NewEncoder(w).Encode(response)
	}
}

// `/todos/{id}` エンドポイントを取得
func TodoHandlers(w http.ResponseWriter, r *http.Request) {
	// リクエスト URL から ID を抽出
	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	// 文字列を整数に変換
	id, err := strconv.Atoi(idStr)
	// 変換に失敗した場合、または ID が空の場合、エラーレスポンスを返す
	if err != nil || idStr == "" {
		http.Error(w, "Invalid ID in URL", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// ストレージ内の指定された ID の Todo を取得
		_, err := handlerStorage.GetByID(id)
		// 取得に失敗した場合、 HTTP 404 not found を返す
		if err != nil {
			http.Error(w, "Todo not found", http.StatusNotFound)
		}
		// レスポンスヘッダに JSON コンテンツタイプを設定
		w.Header().Set("Content-Type", "application/json")
		// 成功メッセージを JSON 形式で返す
		response := map[string]string{
			"message": "Todo fetched successfully",
		}
		// 取得された Todo を JSON 形式でエンコードし、レスポンスとして送信（元のデータ構造 -> JSON 形式に変換）
		json.NewEncoder(w).Encode(response)

	case http.MethodPut:
		// 指定された ID の Todo を更新
		var updateTodo models.Todo
		// リクエストボディから JSON データをデコードして、updateTodo 構造体に格納（JSON 形式 -> 元のデータ構造に変換）
		err := json.NewDecoder(r.Body).Decode(&updateTodo)
		// デコードエラー時、クライアントに400 Bad Requestを返す
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		// 更新する Todo に URL から取得した ID を設定
		updateTodo.ID = id
		// ストレージ内の指定された ID の Todo を新しいデータで更新
		err = handlerStorage.Update(id, updateTodo)
		// 削除に失敗した場合、 HTTP 404 not found を返す
		if err != nil {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		// バリデーションチェック
		if updateTodo.Title == "" || updateTodo.Description == "" {
			http.Error(w, "Please enter your information", http.StatusBadRequest)
			return
		}
		// レスポンスヘッダに JSON コンテンツタイプを設定
		w.Header().Set("Content-Type", "application/json")
		// 成功メッセージを JSON 形式で返す
		response := map[string]string{
			"message": "Todo updated successfully",
		}
		// 更新された Todo を JSON 形式でエンコードし、レスポンスとして送信（元のデータ構造 -> JSON 形式に変換）
		json.NewEncoder(w).Encode(response)

	case http.MethodDelete:
		// 指定された ID の Todo を削除
		err := handlerStorage.Delete(id)
		// 削除に失敗した場合、 HTTP 404 not found を返す
		if err != nil {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		// レスポンスヘッダに JSON コンテンツタイプを設定
		w.Header().Set("Content-Type", "application/json")
		// 削除が成功した場合、 HTTP 204 No Content を返す
		response := map[string]string{
			"message": "Todo deleted successfully",
		}
		// 削除された Todo を JSON 形式でエンコードし、レスポンスとして送信（元のデータ構造 -> JSON 形式に変換）
		json.NewEncoder(w).Encode(response)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
