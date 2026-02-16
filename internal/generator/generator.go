package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Описываем минимально необходимую структуру OpenAPI
type OpenAPI struct {
	Paths map[string]map[string]Operation `yaml:"paths"`
}

type Operation struct {
	OperationID string `yaml:"operationId"`
}

func main() {
	// 1. Читаем файл
	data, err := os.ReadFile("api/api.yml")
	if err != nil {
		log.Fatalf("Ошибка чтения файла: %v", err)
	}

	// 2. Десериализуем YAML
	var spec OpenAPI
	err = yaml.Unmarshal(data, &spec)
	if err != nil {
		log.Fatalf("Ошибка парсинга YAML: %v", err)
	}

	fmt.Println("Найденные OperationID:")
	fmt.Println("----------------------")

	// 3. Перебираем пути и методы
	for path, methods := range spec.Paths {
		for method, details := range methods {
			if details.OperationID != "" {
				fmt.Printf("[%s] %-15s -> %s\n", method, path, details.OperationID)
			}
		}
	}
}
