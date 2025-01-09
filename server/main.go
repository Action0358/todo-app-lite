package main

import (
	"log"
	"net/http"

	"github.com/Action0358/todo-app-lite/server/handlers"
	"github.com/Action0358/todo-app-lite/server/sqlite"
)

func main() {
	// SQLite ストレージの初期化
	storageInstance, err := sqlite.NewSQLiteStorage("")
	if err != nil {
		log.Printf("Faild to initialize SQLite storage: %v", err)
	}

	// ハンドラーにストレージを渡す
	handlers.SetStorage(storageInstance)

	// Restfulなルーティング設定
	http.HandleFunc("/todos", handlers.TodosHandlers) // 一覧取得、新規作成
	http.HandleFunc("/todos/", handlers.TodoHandlers) // 特定のタスクの操作

	// サーバー起動
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("Error starting server: %v", err)
		// エラー発生時に main 関数の処理を終了させる
		return
	}
}
