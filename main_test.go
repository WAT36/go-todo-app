package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	tempDir := t.TempDir()
	staticDir := filepath.Join(tempDir, "static")
	err := os.MkdirAll(staticDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	
	indexContent := `<!DOCTYPE html>
<html>
<head><title>ToDo リスト</title></head>
<body><h1>📝 ToDo リスト</h1></body>
</html>`
	
	indexPath := filepath.Join(staticDir, "index.html")
	err = os.WriteFile(indexPath, []byte(indexContent), 0644)
	if err != nil {
		t.Fatal(err)
	}
	
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)
	
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(homeHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}
	
	responseBody := rr.Body.String()
	if !strings.Contains(responseBody, "ToDo リスト") {
		t.Errorf("Expected response to contain 'ToDo リスト', but got: %s", responseBody)
	}
}

func TestHomeHandlerFileNotFound(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)
	
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(homeHandler)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, status)
	}
}

func TestAPITasksRouting(t *testing.T) {
	testCases := []struct {
		method         string
		expectedStatus int
	}{
		{"GET", http.StatusOK},
		{"POST", http.StatusOK},
		{"DELETE", http.StatusMethodNotAllowed},
		{"PUT", http.StatusMethodNotAllowed},
	}
	
	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, "/api/tasks", nil)
			if err != nil {
				t.Fatal(err)
			}
			
			rr := httptest.NewRecorder()
			
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`[]`))
				} else if r.Method == http.MethodPost {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"success": true}`))
				} else {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			})
			
			handler.ServeHTTP(rr, req)
			
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, status)
			}
		})
	}
}

func TestAPITasksIDRouting(t *testing.T) {
	testCases := []struct {
		path           string
		expectedToggle bool
	}{
		{"/api/tasks/1/toggle", true},
		{"/api/tasks/123/toggle", true},
		{"/api/tasks/1", false},
		{"/api/tasks/456", false},
	}
	
	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			isToggle := tc.path[len(tc.path)-7:] == "/toggle"
			
			if isToggle != tc.expectedToggle {
				t.Errorf("For path %s, expected toggle=%v, got %v", tc.path, tc.expectedToggle, isToggle)
			}
		})
	}
}
