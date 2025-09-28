package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"todo-app/models"
)

func setupTestApp() {
	todoApp = models.NewTodoApp()
}

func TestGetTasksHandler(t *testing.T) {
	setupTestApp()
	
	req, err := http.NewRequest("GET", "/api/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTasksHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}
	
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}
	
	var tasks []models.Task
	err = json.Unmarshal(rr.Body.Bytes(), &tasks)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}
	
	if len(tasks) != 0 {
		t.Errorf("Expected empty tasks array, got %d tasks", len(tasks))
	}
}

func TestGetTasksHandlerWithTasks(t *testing.T) {
	setupTestApp()
	
	todoApp.AddTask("Task 1")
	todoApp.AddTask("Task 2")
	
	req, err := http.NewRequest("GET", "/api/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTasksHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}
	
	var tasks []models.Task
	err = json.Unmarshal(rr.Body.Bytes(), &tasks)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}
	
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
	
	if tasks[0].Title != "Task 1" {
		t.Errorf("Expected first task title 'Task 1', got '%s'", tasks[0].Title)
	}
}

func TestGetTasksHandlerInvalidMethod(t *testing.T) {
	setupTestApp()
	
	req, err := http.NewRequest("POST", "/api/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTasksHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, status)
	}
}

func TestAddTaskHandler(t *testing.T) {
	setupTestApp()
	
	requestBody := map[string]string{
		"title": "New Task",
	}
	jsonBody, _ := json.Marshal(requestBody)
	
	req, err := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddTaskHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}
	
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}
	
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}
	
	if success, ok := response["success"].(bool); !ok || !success {
		t.Error("Expected success to be true")
	}
	
	if task, ok := response["task"].(map[string]interface{}); ok {
		if title, ok := task["title"].(string); !ok || title != "New Task" {
			t.Errorf("Expected task title 'New Task', got '%v'", title)
		}
		if id, ok := task["id"].(float64); !ok || id != 1 {
			t.Errorf("Expected task ID 1, got %v", id)
		}
		if completed, ok := task["completed"].(bool); !ok || completed {
			t.Errorf("Expected task completed false, got %v", completed)
		}
	} else {
		t.Error("Expected task object in response")
	}
}

func TestAddTaskHandlerInvalidMethod(t *testing.T) {
	setupTestApp()
	
	req, err := http.NewRequest("GET", "/api/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddTaskHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, status)
	}
}

func TestAddTaskHandlerInvalidJSON(t *testing.T) {
	setupTestApp()
	
	req, err := http.NewRequest("POST", "/api/tasks", strings.NewReader("invalid json"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddTaskHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
	}
}

func TestAddTaskHandlerEmptyTitle(t *testing.T) {
	setupTestApp()
	
	requestBody := map[string]string{
		"title": "",
	}
	jsonBody, _ := json.Marshal(requestBody)
	
	req, err := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddTaskHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
	}
}

func TestToggleTaskHandler(t *testing.T) {
	setupTestApp()
	
	task := todoApp.AddTask("Test Task")
	
	req, err := http.NewRequest("PUT", "/api/tasks/1/toggle", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ToggleTaskHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}
	
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}
	
	var response map[string]bool
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}
	
	if success, ok := response["success"]; !ok || !success {
		t.Error("Expected success to be true")
	}
	
	tasks := todoApp.GetTasks()
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
	if tasks[0].ID == task.ID && !tasks[0].Completed {
		t.Error("Expected task to be completed after toggle")
	}
}

func TestToggleTaskHandlerInvalidMethod(t *testing.T) {
	setupTestApp()
	
	req, err := http.NewRequest("GET", "/api/tasks/1/toggle", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ToggleTaskHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, status)
	}
}

func TestToggleTaskHandlerInvalidID(t *testing.T) {
	setupTestApp()
	
	req, err := http.NewRequest("PUT", "/api/tasks/invalid/toggle", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ToggleTaskHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
	}
}

func TestToggleTaskHandlerNonExistentTask(t *testing.T) {
	setupTestApp()
	
	req, err := http.NewRequest("PUT", "/api/tasks/999/toggle", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ToggleTaskHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}
	
	var response map[string]bool
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}
	
	if success, ok := response["success"]; !ok || success {
		t.Error("Expected success to be false for non-existent task")
	}
}

func TestDeleteTaskHandler(t *testing.T) {
	setupTestApp()
	
	todoApp.AddTask("Test Task")
	
	req, err := http.NewRequest("DELETE", "/api/tasks/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteTaskHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}
	
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}
	
	var response map[string]bool
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}
	
	if success, ok := response["success"]; !ok || !success {
		t.Error("Expected success to be true")
	}
	
	tasks := todoApp.GetTasks()
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks after deletion, got %d", len(tasks))
	}
}

func TestDeleteTaskHandlerInvalidMethod(t *testing.T) {
	setupTestApp()
	
	req, err := http.NewRequest("GET", "/api/tasks/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteTaskHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, status)
	}
}

func TestDeleteTaskHandlerInvalidID(t *testing.T) {
	setupTestApp()
	
	req, err := http.NewRequest("DELETE", "/api/tasks/invalid", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteTaskHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
	}
}

func TestDeleteTaskHandlerNonExistentTask(t *testing.T) {
	setupTestApp()
	
	req, err := http.NewRequest("DELETE", "/api/tasks/999", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteTaskHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}
	
	var response map[string]bool
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}
	
	if success, ok := response["success"]; !ok || success {
		t.Error("Expected success to be false for non-existent task")
	}
}
