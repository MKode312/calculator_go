package application

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalcHandler(t *testing.T) {
	tests := []struct {
		name           string
		expression     string
		expectedResult string
		expectedStatus int
	}{
		{
			name:           "Valid expression",
			expression:     "2 + 2",
			expectedResult: "result: 4.000000",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid expression",
			expression:     "2 +",
			expectedResult: "expression cannot end with an operator",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Division by zero",
			expression:     "1 / 0",
			expectedResult: "division by zero",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty expression",
			expression:     "",
			expectedResult: "empty expression",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(map[string]string{"expression": tt.expression})
			req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(CalcHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if !bytes.Contains(rr.Body.Bytes(), []byte(tt.expectedResult)) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedResult, rr.Body.String())
			}
		})
	}
}
