package main

import (
	"go/ast"
)

func getProperty(importMap map[string]string, expr ast.Expr) *Property {
	prop := analyzeType(expr)
	if prop == nil {
		return nil
	}
	if prop.IsExported {
		if prop.Alias == "" {
			prop.ImportPath = importMap[prop.Alias]
		} else {
			prop.ImportPath = importMap[prop.Alias]
		}
	}
	return prop
}

func analyzeType(expr ast.Expr) *Property {
	switch t := expr.(type) {
	case *ast.Ident:
		// Локальный тип: ReservationService
		return &Property{
			Name:       t.Name,
			IsExported: ast.IsExported(t.Name),
		}

	case *ast.SelectorExpr:
		// Тип из пакета: repositories.Port
		pkgAlias := ""
		if xIdent, ok := t.X.(*ast.Ident); ok {
			pkgAlias = xIdent.Name
		}
		return &Property{
			Name:       t.Sel.Name,
			Alias:      pkgAlias,
			IsExported: ast.IsExported(t.Sel.Name),
		}

	case *ast.StarExpr:
		// Указатель: *Something
		p := analyzeType(t.X)
		if p != nil {
			p.Name = "*" + p.Name
		}
		return p

	case *ast.ArrayType:
		// Слайс: []Something
		p := analyzeType(t.Elt)
		if p != nil {
			p.Name = "[]" + p.Name
		}
		return p

	default:
		return nil
	}
}
