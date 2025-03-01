package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/MKode312/calculator_go/internal/application"
)

var originalComputingPower = os.Getenv("COMPUTING_POWER")
var originalOrchestratorURL = os.Getenv("ORCHESTRATOR_URL")

func TestNewAgent_DefaultSettings(t *testing.T) {
	os.Unsetenv("COMPUTING_POWER")
	os.Unsetenv("ORCHESTRATOR_URL")

	agent := application.NewAgent()

	if agent.ComputingPower != 1 {
		t.Errorf("expected ComputingPower to be 1, got %d", agent.ComputingPower)
	}
	if agent.OrchestratorURL != "http://localhost:8080" {
		t.Errorf("expected OrchestratorURL to be 'http://localhost:8080', got '%s'", agent.OrchestratorURL)
	}
}

func TestNewAgent_WithEnvironmentVariables(t *testing.T) {
	os.Setenv("COMPUTING_POWER", "3")
	os.Setenv("ORCHESTRATOR_URL", "http://test-url")

	agent := application.NewAgent()

	if agent.ComputingPower != 3 {
		t.Errorf("expected ComputingPower to be 3, got %d", agent.ComputingPower)
	}
	if agent.OrchestratorURL != "http://test-url" {
		t.Errorf("expected OrchestratorURL to be 'http://test-url', got '%s'", agent.OrchestratorURL)
	}

	os.Setenv("COMPUTING_POWER", originalComputingPower)
	os.Setenv("ORCHESTRATOR_URL", originalOrchestratorURL)
}

func TestWorker(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/internal/task" {
			task := map[string]interface{}{
				"task": map[string]interface{}{
					"id":             "task-1",
					"arg1":           5.0,
					"arg2":           3.0,
					"operation":      "add",
					"operation_time": 100,
				},
			}
			json.NewEncoder(w).Encode(task)
		} else if r.Method == http.MethodPost && r.URL.Path == "/internal/task" {
			w.WriteHeader(http.StatusOK)
		}
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	agent := &application.Agent{
		ComputingPower:  1,
		OrchestratorURL: server.URL,
	}
	go agent.Run()

	time.Sleep(300 * time.Millisecond)
}
