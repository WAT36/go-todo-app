package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"todo-app/handlers"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("static", "index.html"))
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	
	http.HandleFunc("/", homeHandler)
	
	http.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetTasksHandler(w, r)
		} else if r.Method == http.MethodPost {
			handlers.AddTaskHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	
	http.HandleFunc("/api/tasks/", func(w http.ResponseWriter, r *http.Request) {
		// /api/tasks/{id}/toggle か /api/tasks/{id} (DELETE) を振り分け
		if r.URL.Path[len(r.URL.Path)-7:] == "/toggle" {
			handlers.ToggleTaskHandler(w, r)
		} else {
			handlers.DeleteTaskHandler(w, r)
		}
	})

	port := "8080"
	fmt.Printf("ToDo アプリケーションを開始しています...\n")
	fmt.Printf("ブラウザで http://localhost:%s にアクセスしてください\n", port)

	// 指定ポートでHTTPサーバを起動（Ctrl+Cで停止）
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
