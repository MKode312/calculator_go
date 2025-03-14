package application

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/MKode312/calculator_go/pkg/ast"
	"github.com/MKode312/calculator_go/pkg/errorsForCalc"
)

type Config struct {
	Addr                string
	TimeAddition        int
	TimeSubtraction     int
	TimeMultiplications int
	TimeDivisions       int
}

func ConfigFromEnv() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	ta, _ := strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
	if ta == 0 {
		ta = 100
	}
	ts, _ := strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
	if ts == 0 {
		ts = 100
	}
	tm, _ := strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
	if tm == 0 {
		tm = 100
	}
	td, _ := strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))
	if td == 0 {
		td = 100
	}
	return &Config{
		Addr:                port,
		TimeAddition:        ta,
		TimeSubtraction:     ts,
		TimeMultiplications: tm,
		TimeDivisions:       td,
	}
}

type Orchestrator struct {
	Config      *Config
	ExprStore   map[string]*Expression
	TaskStore   map[string]*Task
	TaskQueue   []*Task
	mu          sync.Mutex
	ExprCounter int64
	TaskCounter int64
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		Config:    ConfigFromEnv(),
		ExprStore: make(map[string]*Expression),
		TaskStore: make(map[string]*Task),
		TaskQueue: make([]*Task, 0),
	}
}

type Expression struct {
	ID     string       `json:"id"`
	Expr   string       `json:"expression"`
	Status string       `json:"status"`
	Result *float64     `json:"result,omitempty"`
	AST    *ast.ASTNode `json:"-"`
}

type Task struct {
	ID            string       `json:"id"`
	ExprID        string       `json:"-"`
	Arg1          float64      `json:"arg1"`
	Arg2          float64      `json:"arg2"`
	Operation     string       `json:"operation"`
	OperationTime int          `json:"operation_time"`
	Node          *ast.ASTNode `json:"-"`
}

func CheckExpressionErrors(node *ast.ASTNode) error {
	if node == nil {
		return errorsForCalc.ErrEmptyExpression
	}

	if err := ValidateNode(node); err != nil {
		return err
	}

	return nil
}

func ValidateNode(node *ast.ASTNode) error {
	if node.IsLeaf {
		return nil
	}

	// Проверка операторов
	if node.Operator != "+" && node.Operator != "-" && node.Operator != "*" && node.Operator != "/" {
		return errorsForCalc.ErrInvalidOp
	}

	if node.Left == nil || node.Right == nil {
		return fmt.Errorf("missing operand for operator %s", node.Operator)
	}

	if node.Left != nil && !node.Left.IsLeaf {
		if node.Left.Operator != "" {
			return errorsForCalc.ErrTwoConsecutiveOps
		}
	}
	if node.Right != nil && !node.Right.IsLeaf {
		if node.Right.Operator != "" {
			return errorsForCalc.ErrTwoConsecutiveOps
		}
	}

	if node.Operator == "/" && node.Right != nil && node.Right.Value == 0 {
		return errorsForCalc.ErrDivisionByZero
	}

	if err := ValidateNode(node.Left); err != nil {
		return err
	}
	if err := ValidateNode(node.Right); err != nil {
		return err
	}

	return nil
}

func (o *Orchestrator) CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Wrong Method"}`, http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Expression string `json:"expression"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Expression == "" {
		http.Error(w, `{"error":"Invalid Body"}`, http.StatusUnprocessableEntity)
		return
	}
	ast, err := ast.ParseAST(req.Expression)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusUnprocessableEntity)
		return
	}
	o.mu.Lock()
	o.ExprCounter++
	exprID := fmt.Sprintf("%d", o.ExprCounter)
	expr := &Expression{
		ID:     exprID,
		Expr:   req.Expression,
		Status: "pending",
		AST:    ast,
	}
	o.ExprStore[exprID] = expr
	o.ScheduleTasks(expr)
	o.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": exprID})
}

func (o *Orchestrator) ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Wrong Method"}`, http.StatusMethodNotAllowed)
		return
	}
	o.mu.Lock()
	defer o.mu.Unlock()

	exprs := make([]*Expression, 0, len(o.ExprStore))
	for _, expr := range o.ExprStore {
		if expr.AST != nil && expr.AST.IsLeaf {
			if err := CheckExpressionErrors(expr.AST); err != nil {
				expr.Status = "error"
				expr.Result = nil
			} else {
				expr.Status = "completed"
				expr.Result = &expr.AST.Value
			}
		}
		exprs = append(exprs, expr)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"expressions": exprs})
}

func (o *Orchestrator) ExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Wrong Method"}`, http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Path[len("/api/v1/expressions/"):]
	o.mu.Lock()
	expr, ok := o.ExprStore[id]
	o.mu.Unlock()
	if !ok {
		http.Error(w, `{"error":"Expression not found"}`, http.StatusNotFound)
		return
	}
	if expr.AST != nil && expr.AST.IsLeaf {
		expr.Status = "completed"
		expr.Result = &expr.AST.Value
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"expression": expr})
}

func (o *Orchestrator) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Wrong Method"}`, http.StatusMethodNotAllowed)
		return
	}
	o.mu.Lock()
	defer o.mu.Unlock()

	if len(o.TaskQueue) == 0 {
		http.Error(w, `{"error":"No task available"}`, http.StatusNotFound)
		return
	}

	task := o.TaskQueue[0]
	o.TaskQueue = o.TaskQueue[1:]

	if expr, exists := o.ExprStore[task.ExprID]; exists {
		if expr.Status == "pending" || expr.Status == "completed" {
			if err := CheckExpressionErrors(expr.AST); err != nil {
				expr.Status = "error"
				expr.Result = nil
			} else {
				expr.Status = "in_progress"
			}
		}
	} else {
		http.Error(w, `{"error":"Task expression not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"task": task})
}
func (o *Orchestrator) PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Wrong Method"}`, http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ID     string  `json:"id"`
		Result float64 `json:"result"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.ID == "" {
		http.Error(w, `{"error":"Invalid Body"}`, http.StatusUnprocessableEntity)
		return
	}
	o.mu.Lock()
	task, ok := o.TaskStore[req.ID]
	if !ok {
		o.mu.Unlock()
		http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
		return
	}
	task.Node.IsLeaf = true
	task.Node.Value = req.Result
	delete(o.TaskStore, req.ID)
	if expr, exists := o.ExprStore[task.ExprID]; exists {
		o.ScheduleTasks(expr)
		if expr.AST.IsLeaf {
			expr.Status = "completed"
			expr.Result = &expr.AST.Value
		}
	}
	o.mu.Unlock()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"result accepted"}`))
}

func (o *Orchestrator) ScheduleTasks(expr *Expression) {
	var traverse func(node *ast.ASTNode)
	traverse = func(node *ast.ASTNode) {
		if node == nil || node.IsLeaf {
			return
		}
		traverse(node.Left)
		traverse(node.Right)
		if node.Left != nil && node.Right != nil && node.Left.IsLeaf && node.Right.IsLeaf {
			if !node.TaskScheduled {
				o.TaskCounter++
				taskID := fmt.Sprintf("%d", o.TaskCounter)
				var opTime int
				switch node.Operator {
				case "+":
					opTime = o.Config.TimeAddition
				case "-":
					opTime = o.Config.TimeSubtraction
				case "*":
					opTime = o.Config.TimeMultiplications
				case "/":
					opTime = o.Config.TimeDivisions
				default:
					opTime = 100
				}
				task := &Task{
					ID:            taskID,
					ExprID:        expr.ID,
					Arg1:          node.Left.Value,
					Arg2:          node.Right.Value,
					Operation:     node.Operator,
					OperationTime: opTime,
					Node:          node,
				}
				node.TaskScheduled = true
				o.TaskStore[taskID] = task
				o.TaskQueue = append(o.TaskQueue, task)
			}
		}
	}
	traverse(expr.AST)
}

func (o *Orchestrator) RunServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", o.CalculateHandler)
	mux.HandleFunc("/api/v1/expressions", o.ExpressionsHandler)
	mux.HandleFunc("/api/v1/expressions/", o.ExpressionByIDHandler)
	mux.HandleFunc("/internal/task", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			o.GetTaskHandler(w, r)
		} else if r.Method == http.MethodPost {
			o.PostTaskHandler(w, r)
		} else {
			http.Error(w, `{"error":"Wrong Method"}`, http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"Not Found"}`, http.StatusNotFound)
	})
	go func() {
		for {
			time.Sleep(2 * time.Second)
			o.mu.Lock()
			if len(o.TaskQueue) > 0 {
				log.Printf("Pending tasks in queue: %d", len(o.TaskQueue))
			}
			o.mu.Unlock()
		}
	}()
	return http.ListenAndServe(":"+o.Config.Addr, mux)
}
