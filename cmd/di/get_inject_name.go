package main

import (
	"go/ast"
	"regexp"
)

var injectRegexp = regexp.MustCompile(`di:use\(([^)]+)\)`)

func getInjectName(field *ast.Field, commentMap ast.CommentMap) string {
	comments := commentMap[field]
	for _, comment := range comments {
		matches := injectRegexp.FindStringSubmatch(comment.Text())
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}
