package names

import (
	"regexp"
	"strings"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func ToSnakeCase(str string) string {
	// 1. Вставляем подчеркивание перед заглавными буквами, за которыми следуют строчные
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	// 2. Вставляем подчеркивание между строчной (или цифрой) и заглавной
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	// 3. Все в нижний регистр
	return strings.ToLower(snake)
}
