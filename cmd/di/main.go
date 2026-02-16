package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	root := "."
	//root = "/Users/emris/www/gor/example/internal"
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	outputDir := filepath.Join(root, "generated", "di")
	if err := os.RemoveAll(outputDir); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to clean di directory - %v\n", err)
		os.Exit(1)
	}

	dependencies, err := scanDependencies(root)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "DI: failed - %v\n", err)
		os.Exit(1)
	}
	//data, _ := json.MarshalIndent(dependencies, "", "    ")
	//fmt.Println(string(data))

	err = generate(outputDir, dependencies)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "DI: failed write files - %v\n", err)
		os.Exit(1)
	}
}
