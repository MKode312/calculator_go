package application

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/MKode312/calculator_go/pkg/calc"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

type CalcRequest struct {
	Exp string `json:"expression"`
}

type CalcResponse struct {
	Result float64 `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, errorCode int, errorText string) {
	errorResponse := ErrorResponse{
		Error: errorText,
	}
	data, err := json.Marshal(errorResponse)
	if err != nil {
		_ = fmt.Errorf("failed marshal response: %w", err)
		return
	}
	w.WriteHeader(errorCode)
	_, err = w.Write(data)
	if err != nil {
		_ = fmt.Errorf("failed to write response: %w", err)
	}
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println("Failed to close request body")
		}
	}(r.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	var calcRequest CalcRequest
	err = json.Unmarshal(body, &calcRequest)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "Expression is not valid")
		return
	}

	ans, err := calc.Calc(calcRequest.Exp)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "Expression is not valid")
		return
	}

	calculateResponse := CalcResponse{
		Result: ans,
	}

	data, err := json.Marshal(calculateResponse)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
	} else {
		_, err = w.Write(data)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Internal server error")
		}
	}
}

func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
