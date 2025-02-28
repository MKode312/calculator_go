package main

import (
	"log"

	"github.com/MKode312/calculator_go/internal/application"
)

func main() {
	agent := application.NewAgent()
	log.Println("Starting Agent...")
	agent.Run()
}
