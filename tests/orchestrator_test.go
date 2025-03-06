package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/MKode312/calculator_go/internal/application"
	"github.com/MKode312/calculator_go/pkg/ast"
)

// Test ConfigFromEnv function
func TestConfigFromEnv(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("PORT", "8081")
	os.Setenv("TIME_ADDITION_MS", "200")
	os.Setenv("TIME_SUBTRACTION_MS", "200")
	os.Setenv("TIME_MULTIPLICATIONS_MS", "300")
	os.Setenv("TIME_DIVISIONS_MS", "400")
	defer os.Clearenv()

	config := application.ConfigFromEnv()

	if config.Addr != "8081" {
		t.Errorf("Expected PORT to be 8081, got %s", config.Addr)
	}
	if config.TimeAddition != 200 {
		t.Errorf("Expected TIME_ADDITION_MS to be 200, got %d", config.TimeAddition)
	}
	if config.TimeSubtraction != 200 {
		t.Errorf("Expected TIME_SUBTRACTION_MS to be 200, got %d", config.TimeSubtraction)
	}
	if config.TimeMultiplications != 300 {
		t.Errorf("Expected TIME_MULTIPLICATIONS_MS to be 300, got %d", config.TimeMultiplications)
	}
	if config.TimeDivisions != 400 {
		t.Errorf("Expected TIME_DIVISIONS_MS to be 400, got %d", config.TimeDivisions)
	}
}

// Test Orchestrator's calculateHandler method
func TestCalculateHandler(t *testing.T) {
	orchestrator := application.NewOrchestrator()
	// Добавить выражение или задачу перед тестом, чтобы использовать правильный ID
	orchestrator.ExprStore["1"] = &application.Expression{
		ID:     "1",
		Expr:   "3 + 2",
		Status: "pending",
		AST:    &ast.ASTNode{IsLeaf: false},
	}

	go func() {
		err := orchestrator.RunServer()
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("Server exited with error: %v", err)
		}
	}()

}
func TestExpressionsHandler(t *testing.T) {
	orchestrator := application.NewOrchestrator()
	orchestrator.ExprStore["1"] = &application.Expression{
		ID:     "1",
		Expr:   "2 + 2",
		Status: "completed",
		AST:    &ast.ASTNode{IsLeaf: true, Value: 4},
	}

	tests := []struct {
		method      string
		statusCode  int
		expectedLen int // Expected length of the expressions array
	}{
		{http.MethodGet, http.StatusOK, 1},
		{"POST", http.StatusMethodNotAllowed, 0},
	}

	for _, test := range tests {
		req, _ := http.NewRequest(test.method, "/api/v1/expressions", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(orchestrator.ExpressionsHandler)

		handler.ServeHTTP(rr, req)

		if rr.Code != test.statusCode {
			t.Errorf("For method %s, expected status code %d but got %d", test.method, test.statusCode, rr.Code)
		}

		if test.statusCode == http.StatusOK {
			var response map[string]interface{}
			json.Unmarshal(rr.Body.Bytes(), &response)
			if len(response["expressions"].([]interface{})) != test.expectedLen {
				t.Errorf("Expected %d expressions, got %d", test.expectedLen, len(response["expressions"].([]interface{})))
			}
		}
	}
}

func TestExpressionByIDHandler(t *testing.T) {
	orchestrator := application.NewOrchestrator()
	orchestrator.ExprStore["1"] = &application.Expression{
		ID:     "1",
		Expr:   "2 + 2",
		Status: "completed",
		AST:    &ast.ASTNode{IsLeaf: true, Value: 4},
	}

	tests := []struct {
		id         string
		method     string
		statusCode int
	}{
		{"1", http.MethodGet, http.StatusOK},
		{"2", http.MethodGet, http.StatusNotFound},
		{"1", "POST", http.StatusMethodNotAllowed},
	}

	for _, test := range tests {
		req, _ := http.NewRequest(test.method, "/api/v1/expressions/"+test.id, nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(orchestrator.ExpressionByIDHandler)

		handler.ServeHTTP(rr, req)

		if rr.Code != test.statusCode {
			t.Errorf("For ID %s with method %s, expected status code %d but got %d", test.id, test.method, test.statusCode, rr.Code)
		}
	}
}

func TestGetTaskHandler(t *testing.T) {
	orchestrator := application.NewOrchestrator()
	orchestrator.TaskQueue = []*application.Task{{ID: "1", ExprID: "1", Arg1: 3, Arg2: 2, Operation: "+"}}
	orchestrator.ExprStore["1"] = &application.Expression{
		ID:     "1",
		Expr:   "3 + 2",
		Status: "pending",
		AST:    &ast.ASTNode{IsLeaf: false},
	}

	tests := []struct {
		method     string
		statusCode int
	}{
		{http.MethodGet, http.StatusOK},
		{"POST", http.StatusMethodNotAllowed},
	}

	for _, test := range tests {
		req, _ := http.NewRequest(test.method, "/api/v1/tasks", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(orchestrator.GetTaskHandler)

		handler.ServeHTTP(rr, req)

		if rr.Code != test.statusCode {
			t.Errorf("For method %s, expected status code %d but got %d", test.method, test.statusCode, rr.Code)
		}

		if test.statusCode == http.StatusOK {
			var response map[string]interface{}
			json.Unmarshal(rr.Body.Bytes(), &response)
			if response["task"] == nil {
				t.Error("Expected task in response, got nil")
			}
		}
	}
}

func TestPostTaskHandler(t *testing.T) {
	orchestrator := application.NewOrchestrator()
	orchestrator.ExprStore["1"] = &application.Expression{
		ID:     "1",
		Expr:   "3 + 2",
		Status: "pending",
		AST:    &ast.ASTNode{IsLeaf: false},
	}

	orchestrator.TaskStore["1"] = &application.Task{ID: "1", ExprID: "1", Node: &ast.ASTNode{}}

	tests := []struct {
		method     string
		body       string
		statusCode int
	}{
		{http.MethodPost, `{"id":"1","result":5}`, http.StatusOK},
		{"GET", "", http.StatusMethodNotAllowed},
		{http.MethodPost, `{"id":"","result":5}`, http.StatusUnprocessableEntity},
		{http.MethodPost, `invalid json`, http.StatusUnprocessableEntity},
		{http.MethodPost, `{"id":"2","result":5}`, http.StatusNotFound},
	}

	for _, test := range tests {
		req, _ := http.NewRequest(test.method, "/api/v1/tasks", bytes.NewBufferString(test.body))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(orchestrator.PostTaskHandler)

		handler.ServeHTTP(rr, req)

		if rr.Code != test.statusCode {
			t.Errorf("For method %s with body %s, expected status code %d but got %d", test.method, test.body, test.statusCode, rr.Code)
		}
	}
}
