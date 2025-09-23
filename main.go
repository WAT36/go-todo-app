package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type TodoApp struct {
	tasks  []Task
	nextID int
	mutex  sync.RWMutex
}

func NewTodoApp() *TodoApp {
	return &TodoApp{
		tasks:  make([]Task, 0),
		nextID: 1,
	}
}

func (app *TodoApp) AddTask(title string) Task {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	
	task := Task{
		ID:        app.nextID,
		Title:     title,
		Completed: false,
	}
	app.tasks = append(app.tasks, task)
	app.nextID++
	return task
}

func (app *TodoApp) GetTasks() []Task {
	app.mutex.RLock()
	defer app.mutex.RUnlock()
	
	tasksCopy := make([]Task, len(app.tasks))
	copy(tasksCopy, app.tasks)
	return tasksCopy
}

func (app *TodoApp) ToggleTask(id int) bool {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	
	for i := range app.tasks {
		if app.tasks[i].ID == id {
			app.tasks[i].Completed = !app.tasks[i].Completed
			return true
		}
	}
	return false
}

func (app *TodoApp) DeleteTask(id int) bool {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	
	for i, task := range app.tasks {
		if task.ID == id {
			app.tasks = append(app.tasks[:i], app.tasks[i+1:]...)
			return true
		}
	}
	return false
}

var todoApp = NewTodoApp()

const htmlTemplate = `
<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ToDo リスト</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background: white;
            border-radius: 10px;
            padding: 30px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            text-align: center;
            margin-bottom: 30px;
        }
        .add-task {
            display: flex;
            gap: 10px;
            margin-bottom: 30px;
        }
        .add-task input {
            flex: 1;
            padding: 12px;
            border: 2px solid #ddd;
            border-radius: 5px;
            font-size: 16px;
        }
        .add-task input:focus {
            outline: none;
            border-color: #4CAF50;
        }
        .add-task button {
            padding: 12px 20px;
            background: #4CAF50;
            color: white;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            font-size: 16px;
        }
        .add-task button:hover {
            background: #45a049;
        }
        .task-list {
            list-style: none;
            padding: 0;
        }
        .task-item {
            display: flex;
            align-items: center;
            padding: 15px;
            margin-bottom: 10px;
            background: #f9f9f9;
            border-radius: 5px;
            border-left: 4px solid #4CAF50;
        }
        .task-item.completed {
            opacity: 0.6;
            border-left-color: #ccc;
        }
        .task-item.completed .task-title {
            text-decoration: line-through;
            color: #888;
        }
        .task-title {
            flex: 1;
            margin: 0 15px;
            font-size: 16px;
        }
        .task-checkbox {
            width: 20px;
            height: 20px;
            cursor: pointer;
        }
        .delete-btn {
            background: #f44336;
            color: white;
            border: none;
            padding: 8px 12px;
            border-radius: 3px;
            cursor: pointer;
            font-size: 14px;
        }
        .delete-btn:hover {
            background: #da190b;
        }
        .empty-state {
            text-align: center;
            color: #888;
            font-style: italic;
            padding: 40px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>📝 ToDo リスト</h1>
        
        <div class="add-task">
            <input type="text" id="taskInput" placeholder="新しいタスクを入力してください..." maxlength="100">
            <button onclick="addTask()">追加</button>
        </div>
        
        <ul class="task-list" id="taskList">
            {{range .}}
            <li class="task-item {{if .Completed}}completed{{end}}">
                <input type="checkbox" class="task-checkbox" {{if .Completed}}checked{{end}} 
                       onchange="toggleTask({{.ID}})">
                <span class="task-title">{{.Title}}</span>
                <button class="delete-btn" onclick="deleteTask({{.ID}})">削除</button>
            </li>
            {{else}}
            <li class="empty-state">タスクがありません。上記のフォームから新しいタスクを追加してください。</li>
            {{end}}
        </ul>
    </div>

    <script>
        function addTask() {
            const input = document.getElementById('taskInput');
            const title = input.value.trim();
            
            if (!title) {
                alert('タスクの内容を入力してください');
                return;
            }
            
            fetch('/api/tasks', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ title: title })
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    input.value = '';
                    location.reload();
                } else {
                    alert('タスクの追加に失敗しました');
                }
            })
            .catch(error => {
                console.error('Error:', error);
                alert('エラーが発生しました');
            });
        }
        
        function toggleTask(id) {
            fetch('/api/tasks/' + id + '/toggle', {
                method: 'PUT'
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    location.reload();
                } else {
                    alert('タスクの更新に失敗しました');
                }
            })
            .catch(error => {
                console.error('Error:', error);
                alert('エラーが発生しました');
            });
        }
        
        function deleteTask(id) {
            if (confirm('このタスクを削除しますか？')) {
                fetch('/api/tasks/' + id, {
                    method: 'DELETE'
                })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        location.reload();
                    } else {
                        alert('タスクの削除に失敗しました');
                    }
                })
                .catch(error => {
                    console.error('Error:', error);
                    alert('エラーが発生しました');
                });
            }
        }
        
        document.getElementById('taskInput').addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                addTask();
            }
        });
    </script>
</body>
</html>
`

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	tasks := todoApp.GetTasks()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, tasks)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		Title string `json:"title"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	
	task := todoApp.AddTask(req.Title)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"task":    task,
	})
}

func toggleTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	idStr := r.URL.Path[len("/api/tasks/"):]
	idStr = idStr[:len(idStr)-len("/toggle")]
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	
	success := todoApp.ToggleTask(id)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{
		"success": success,
	})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	idStr := r.URL.Path[len("/api/tasks/"):]
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	
	success := todoApp.DeleteTask(id)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{
		"success": success,
	})
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/tasks", addTaskHandler)
	http.HandleFunc("/api/tasks/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path[len(r.URL.Path)-7:] == "/toggle" {
			toggleTaskHandler(w, r)
		} else {
			deleteTaskHandler(w, r)
		}
	})
	
	port := "8080"
	fmt.Printf("ToDo アプリケーションを開始しています...\n")
	fmt.Printf("ブラウザで http://localhost:%s にアクセスしてください\n", port)
	
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
