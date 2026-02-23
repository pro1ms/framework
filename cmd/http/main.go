package main

import (
	"fmt"
	"os"

	"github.com/pro1ms/framework/internal/cmd/http"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Gen HTTP: not enough arguments")
		return
	}

	generator := http.NewScanner()
	err := generator.Run(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
