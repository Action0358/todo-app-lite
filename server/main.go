package main

import (
	"log"
	"net/http"

	"github.com/Action0358/todo-app-lite/server/handlers"
)

func main() {
	// Restfulなルーティング設定
	http.HandleFunc("/todos", handlers.TodosHandler)  // 一覧取得、新規作成
	http.HandleFunc("/todos/", handlers.TodoHandler) // 特定のタスクの操作

	// サーバー起動
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Error starting server", err)
		// エラー発生時に main 関数の処理を終了させる
		return
	}
}
