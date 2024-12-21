package application

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/hxzzf/calc_go/pkg/calculation"
)

type Config struct {
	Port string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Port = os.Getenv("PORT")
	if config.Port == "" {
		config.Port = "8080"
	}
	return config
}

type Application struct {
	config *Config
	server *http.Server
}

func New() *Application {
	config := ConfigFromEnv()
	return &Application{
		config: config,
		server: &http.Server{
			Addr: ":" + config.Port,
		},
	}
}

type CalculateRequest struct {
	Expression string `json:"expression"`
}

type CalculateResponse struct {
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

func (a *Application) RunServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", a.HandleCalculate)
	a.server.Handler = mux

	log.Printf("Starting server on port %s\n", a.config.Port)
	return a.server.ListenAndServe()
}

func (a *Application) HandleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var rawReq map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&rawReq); err != nil {
		sendError(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}

	expressionVal, ok := rawReq["expression"]
	if !ok {
		sendError(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}

	expression, ok := expressionVal.(string)
	if !ok {
		sendError(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}

	if expression == "" {
		sendError(w, "Expression cannot be empty", http.StatusUnprocessableEntity)
		return
	}

	result, err := calculation.Calc(expression)
	if err != nil {
		switch err.Error() {
		case "empty expression":
			sendError(w, "Expression cannot be empty", http.StatusUnprocessableEntity)
		case "division by zero":
			sendError(w, "Division by zero is not allowed", http.StatusUnprocessableEntity)
		case "invalid expression: consecutive operators":
			sendError(w, "Consecutive operators are not allowed", http.StatusUnprocessableEntity)
		case "mismatched parentheses":
			sendError(w, "Parentheses are mismatched", http.StatusUnprocessableEntity)
		case "invalid token", "invalid expression":
			sendError(w, "Expression is not valid", http.StatusUnprocessableEntity)
		default:
			sendError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	sendResponse(w, CalculateResponse{Result: result})
}

func sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(CalculateResponse{Error: message})
}

func sendResponse(w http.ResponseWriter, response CalculateResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (a *Application) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
