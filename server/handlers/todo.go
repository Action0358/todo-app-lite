package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Action0358/todo-app-lite/server/models"
	"github.com/Action0358/todo-app-lite/server/storage"
)

// メモリーストレージの新しいインスタンスを作成
var storageInstance = storage.NewMemoryStorage()

// `/todos` エンドポイントを処理
func TodosHandlers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// 全ての Todo を取得
		todos, err := storageInstance.GetAll()
		// 取得に失敗した場合、 HTTP 500 エラーを返す
		if err != nil {
			http.Error(w, "Failed to retrieve todos", http.StatusInternalServerError)
			return
		}
		// レスポンスヘッダに JSON コンテンツタイプを設定
		w.Header().Set("Content-Type", "application/json")
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
		// ストレージに新しい Todo を追加
		err = storageInstance.Add(newTodo)
		// 追跡処理でエラーが発生した場合、 500 Internal Server Error を返す
		if err != nil {
			http.Error(w, "Failed to add todo", http.StatusInternalServerError)
			return
		}
		// レスポンスヘッダに JSON コンテンツタイプを設定
		w.Header().Set("Content-Type", "application/json")
		// 追跡された Todo を JSON で応答(元のデータ構造 -> JSON 形式に変換)
		json.NewEncoder(w).Encode(newTodo)
		return
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
		todo, err := storageInstance.GetByID(id)
		// 取得に失敗した場合、 HTTP 404 not found を返す
		if err != nil {
			http.Error(w, "Todo not found", http.StatusNotFound)
		}
		// レスポンスヘッダに JSON コンテンツタイプを設定
		w.Header().Set("Content-Type", "application/json")
		// 取得した Todo の一覧を JSON 形式でエンコードし、レスポンスとして送信（元のデータ構造 -> JSON 形式に変換）
		json.NewEncoder(w).Encode(todo)

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
		err = storageInstance.Update(id, updateTodo)
		// 削除に失敗した場合、 HTTP 404 not found を返す
		if err != nil {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		// レスポンスヘッダに JSON コンテンツタイプを設定
		w.Header().Set("Content-Type", "application/json")
		// ユーザーがリクエストしたデータ（updateTodo） を JSON で応答（元のデータ構造 -> JSON 形式に変換）
		json.NewEncoder(w).Encode(updateTodo)

	case http.MethodDelete:
		// 指定された ID の Todo を削除
		err := storageInstance.Delete(id)
		// 削除に失敗した場合、 HTTP 404 not found を返す
		if err != nil {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		// 成功時、 HTTP ステータスコード 204 No Content を設定（レスポンスボディなし）
		w.WriteHeader(http.StatusNoContent)

	default:
		// サポートされていない HTTP メソッドがリクエストされた場合、 405 Method Not Allowed を返す
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
