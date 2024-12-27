package application_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MKode312/calculator_go/internal/application"
)

func TestRequestHandlerOk(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/calculate",
		strings.NewReader(`{"expression": "2 + 3"}`))
	w := httptest.NewRecorder()
	application.CalcHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("wrong status code")
	}
}

func TestRequestHandlerInvalidExpression(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/calculate",
		strings.NewReader(`{"expression": "2 / 0"}`))
	w := httptest.NewRecorder()
	application.CalcHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("wrong status code")
	}
}
