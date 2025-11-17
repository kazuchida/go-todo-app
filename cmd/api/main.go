package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"go_todo_app/internal/adapter/handler"
	"go_todo_app/internal/adapter/repository"
	"go_todo_app/internal/infrastructure/database"
	"go_todo_app/internal/usecase"
)

func main() {
	// データベース接続
	db, err := database.NewPostgresDB(database.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   "todoapp",
	})
	if err != nil {
		log.Fatal("データベース接続エラー:", err)
	}
	defer db.Close()

	// スキーマ初期化
	if err := database.InitSchema(db); err != nil {
		log.Fatal("スキーマ初期化エラー:", err)
	}

	// 依存性注入(外側から内側へ)
	todoRepo := repository.NewTodoRepository(db)
	todoUseCase := usecase.NewTodoUseCase(todoRepo)
	todoHandler := handler.NewTodoHandler(todoUseCase)

	// ルーティング (gorilla/mux を使用)
	router := mux.NewRouter()

	// CORS対応ミドルウェア (http.Handler を受け取る形)
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	// NOTE: router.Use によるミドルウェアはマッチするルートがある場合にのみ
	// 実行されます。OPTIONS (プリフライト) のようにメソッドがマッチしない場合
	// ミドルウェアが呼ばれず CORS ヘッダが返らないことがあるため、
	// ルーターを外側からラップして常に CORS を適用するようにします（下で使用）。

	// ルート登録（メソッドとパスパラメータを指定）
	router.HandleFunc("/api/todos", todoHandler.CreateTodo).Methods("POST")
	router.HandleFunc("/api/todos", todoHandler.GetAllTodos).Methods("GET")
	router.HandleFunc("/api/todos/{id}", todoHandler.GetTodo).Methods("GET")
	router.HandleFunc("/api/todos/{id}", todoHandler.UpdateTodo).Methods("PUT")
	router.HandleFunc("/api/todos/{id}/toggle", todoHandler.ToggleTodo).Methods("PATCH")
	router.HandleFunc("/api/todos/{id}", todoHandler.DeleteTodo).Methods("DELETE")

	// ルーターを CORS ミドルウェアでラップして常にヘッダを返す
	handler := corsMiddleware(router)

	log.Println("サーバー起動: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
