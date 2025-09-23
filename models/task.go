package models

import "sync"

// Task は1件のタスク（やること）を表すデータ構造です
// ID: 一意に識別する番号
// Title: タスクの内容
// Completed: 完了しているかどうか
type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// TodoApp はアプリ全体の状態を管理します
// tasks: すべてのタスク一覧
// nextID: 次に採番するID
// mutex: 複数のリクエストから同時に触られても安全にするためのロック
type TodoApp struct {
	tasks  []Task
	nextID int
	mutex  sync.RWMutex
}

// NewTodoApp は TodoApp の初期化（コンストラクタ）を行います
func NewTodoApp() *TodoApp {
	return &TodoApp{
		tasks:  make([]Task, 0),
		nextID: 1,
	}
}

// AddTask は新しいタスクを作成して一覧に追加します
// 排他ロック（書き込み用）を使って安全に配列へ追加します
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

// GetTasks は現在のタスク一覧をコピーして返します
// 読み取り専用ロックを使い、呼び出し側が書き換えても
// 元データに影響しないようスライスのコピーを返します
func (app *TodoApp) GetTasks() []Task {
	app.mutex.RLock()
	defer app.mutex.RUnlock()

	tasksCopy := make([]Task, len(app.tasks))
	copy(tasksCopy, app.tasks)
	return tasksCopy
}

// ToggleTask は指定IDのタスクの完了フラグを反転（true/false）します
// 見つかったら true を、見つからなければ false を返します
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

// DeleteTask は指定IDのタスクを一覧から削除します
// 見つかったら true を、見つからなければ false を返します
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
