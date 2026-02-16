package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
)

func inspectFile(packagePath string, fileSet *token.FileSet, file *ast.File, path string) ([]Dependency, error) {
	var err error
	var deps []Dependency
	var diNewPattern = regexp.MustCompile(`di:new\(([^)]+)\)`)

	importMap := getImportMap(file)

	commentMap := ast.NewCommentMap(fileSet, file, file.Comments)

	ast.Inspect(file, func(n ast.Node) bool {
		if err != nil {
			return false
		}

		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Doc == nil {
			return true
		}

		for _, comment := range fn.Doc.List {
			matches := diNewPattern.FindStringSubmatch(comment.Text)
			if len(matches) > 1 {
				name := matches[1]
				fmt.Printf("DI: found constructor %s in file %s\n", name, path)

				dep := Dependency{
					Name:        name,
					Func:        fn.Name.Name,
					PackagePath: packagePath,
					PackageName: file.Name.Name,
				}

				for _, field := range fn.Type.Params.List {
					prop := getProperty(importMap, field.Type)
					if prop == nil {
						err = fmt.Errorf("di: failed scan property %s in file %s", name, path)
						return false
					}

					injectName := getInjectName(field, commentMap)
					if injectName == "" {
						err = fmt.Errorf("di: not found inject name for property %s in file %s", name, path)
						return false
					}
					prop.Inject = injectName

					for range field.Names {
						dep.Props = append(dep.Props, *prop)
					}
				}

				deps = append(deps, dep)
			}
		}
		return false
	})

	if err != nil {
		return nil, err
	}
	return deps, nil
}
