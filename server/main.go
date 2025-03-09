package main

import (
	"log"
	"net/http"

	"github.com/Action0358/todo-app-lite/server/handlers"
	"github.com/Action0358/todo-app-lite/server/sqlite"
	"github.com/rs/cors"
)

func main() {
	// SQLite ストレージの初期化
	storageInstance, err := sqlite.NewSQLiteStorage("")
	if err != nil {
		log.Fatalf("Faild to initialize SQLite storage: %v", err)
	}
	// 関数終了時に必ずクローズ
	defer storageInstance.DB.Close()

	// ハンドラーにストレージを渡す
	handlers.SetStorage(storageInstance)

	// CORS 設定（クロスオリジン対応）
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// 高度なルーティング
	mux := http.NewServeMux()
	mux.HandleFunc("/todos", handlers.TodosHandlers)
	mux.HandleFunc("/todos/", handlers.TodoHandlers)

	// CORS ミドルウェアでラップしたハンドラー
	handler := corsHandler.Handler(mux)

	// サーバー起動
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
