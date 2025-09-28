package models

import (
	"sync"
	"testing"
)

func TestNewTodoApp(t *testing.T) {
	app := NewTodoApp()
	
	if app == nil {
		t.Fatal("NewTodoApp() returned nil")
	}
	
	if app.nextID != 1 {
		t.Errorf("Expected nextID to be 1, got %d", app.nextID)
	}
	
	if len(app.tasks) != 0 {
		t.Errorf("Expected empty tasks slice, got length %d", len(app.tasks))
	}
}

func TestAddTask(t *testing.T) {
	app := NewTodoApp()
	
	task := app.AddTask("Test task")
	
	if task.ID != 1 {
		t.Errorf("Expected task ID to be 1, got %d", task.ID)
	}
	
	if task.Title != "Test task" {
		t.Errorf("Expected task title to be 'Test task', got '%s'", task.Title)
	}
	
	if task.Completed != false {
		t.Errorf("Expected task to be incomplete, got %v", task.Completed)
	}
	
	task2 := app.AddTask("Second task")
	if task2.ID != 2 {
		t.Errorf("Expected second task ID to be 2, got %d", task2.ID)
	}
	
	tasks := app.GetTasks()
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
}

func TestAddTaskConcurrency(t *testing.T) {
	app := NewTodoApp()
	var wg sync.WaitGroup
	numGoroutines := 100
	
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			app.AddTask("Task " + string(rune(i)))
		}(i)
	}
	
	wg.Wait()
	
	tasks := app.GetTasks()
	if len(tasks) != numGoroutines {
		t.Errorf("Expected %d tasks, got %d", numGoroutines, len(tasks))
	}
	
	idMap := make(map[int]bool)
	for _, task := range tasks {
		if idMap[task.ID] {
			t.Errorf("Duplicate task ID found: %d", task.ID)
		}
		idMap[task.ID] = true
	}
}

func TestGetTasks(t *testing.T) {
	app := NewTodoApp()
	
	tasks := app.GetTasks()
	if len(tasks) != 0 {
		t.Errorf("Expected empty tasks slice, got length %d", len(tasks))
	}
	
	app.AddTask("Task 1")
	app.AddTask("Task 2")
	
	tasks = app.GetTasks()
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
	
	tasks[0].Title = "Modified"
	originalTasks := app.GetTasks()
	if originalTasks[0].Title == "Modified" {
		t.Error("GetTasks() should return a copy, not the original slice")
	}
}

func TestToggleTask(t *testing.T) {
	app := NewTodoApp()
	
	success := app.ToggleTask(999)
	if success {
		t.Error("Expected ToggleTask to return false for non-existent task")
	}
	
	task := app.AddTask("Test task")
	
	success = app.ToggleTask(task.ID)
	if !success {
		t.Error("Expected ToggleTask to return true for existing task")
	}
	
	tasks := app.GetTasks()
	if !tasks[0].Completed {
		t.Error("Expected task to be completed after toggle")
	}
	
	success = app.ToggleTask(task.ID)
	if !success {
		t.Error("Expected ToggleTask to return true for existing task")
	}
	
	tasks = app.GetTasks()
	if tasks[0].Completed {
		t.Error("Expected task to be incomplete after second toggle")
	}
}

func TestToggleTaskConcurrency(t *testing.T) {
	app := NewTodoApp()
	task := app.AddTask("Test task")
	
	var wg sync.WaitGroup
	numGoroutines := 100
	
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			app.ToggleTask(task.ID)
		}()
	}
	
	wg.Wait()
	
	tasks := app.GetTasks()
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
}

func TestDeleteTask(t *testing.T) {
	app := NewTodoApp()
	
	success := app.DeleteTask(999)
	if success {
		t.Error("Expected DeleteTask to return false for non-existent task")
	}
	
	task1 := app.AddTask("Task 1")
	task2 := app.AddTask("Task 2")
	task3 := app.AddTask("Task 3")
	
	success = app.DeleteTask(task2.ID)
	if !success {
		t.Error("Expected DeleteTask to return true for existing task")
	}
	
	tasks := app.GetTasks()
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks after deletion, got %d", len(tasks))
	}
	
	foundTask2 := false
	for _, task := range tasks {
		if task.ID == task2.ID {
			foundTask2 = true
		}
	}
	if foundTask2 {
		t.Error("Task 2 should have been deleted")
	}
	
	foundTask1 := false
	foundTask3 := false
	for _, task := range tasks {
		if task.ID == task1.ID {
			foundTask1 = true
		}
		if task.ID == task3.ID {
			foundTask3 = true
		}
	}
	if !foundTask1 || !foundTask3 {
		t.Error("Task 1 and Task 3 should still exist")
	}
}

func TestDeleteTaskConcurrency(t *testing.T) {
	app := NewTodoApp()
	
	for i := 0; i < 10; i++ {
		app.AddTask("Task " + string(rune(i)))
	}
	
	var wg sync.WaitGroup
	numGoroutines := 5
	
	wg.Add(numGoroutines)
	for i := 1; i <= numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			app.DeleteTask(id)
		}(i)
	}
	
	wg.Wait()
	
	tasks := app.GetTasks()
	if len(tasks) != 5 {
		t.Errorf("Expected 5 tasks remaining, got %d", len(tasks))
	}
}

func TestTaskStruct(t *testing.T) {
	task := Task{
		ID:        1,
		Title:     "Test Task",
		Completed: true,
	}
	
	if task.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", task.ID)
	}
	
	if task.Title != "Test Task" {
		t.Errorf("Expected Title to be 'Test Task', got '%s'", task.Title)
	}
	
	if task.Completed != true {
		t.Errorf("Expected Completed to be true, got %v", task.Completed)
	}
}
