package main

import (
	"fmt"
	// . "pawwander/table_defs"
	"regexp"
	"strings"
)

func convertToSnakeCase(name string) string {
	// 使用正则表达式找到大写字符并在前面加上下划线，然后转换为小写
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(name, "${1}_${2}")
	return strings.ToLower(snake)
}

func main() {
	fmt.Println(convertToSnakeCase("UserActivity"))
}
