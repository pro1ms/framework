package main

import (
	"go/ast"
	"strings"
)

func getImportMap(file *ast.File) map[string]string {
	res := make(map[string]string)
	for _, imp := range file.Imports {
		// Убираем кавычки: "://github.com..." -> ://github.com...
		path := strings.Trim(imp.Path.Value, `"`)

		var alias string
		if imp.Name != nil {
			alias = imp.Name.Name // Явный алиас: import api "..."
		} else {
			// Алиас по умолчанию (последний сегмент пути)
			parts := strings.Split(path, "/")
			alias = parts[len(parts)-1]
		}
		res[alias] = path
	}
	return res
}
