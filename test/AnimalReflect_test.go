package test

import (
	model "catface/app/model"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unicode"
)

// 获取处理后的 json 标签值
func getProcessedJSONTag(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag != "" {
		// 移除 omitempty
		jsonTag = strings.Replace(jsonTag, ",omitempty", "", -1)
		jsonTag = strings.Replace(jsonTag, "omitempty", "", -1)
	}
	return jsonTag
}

func Test(tt *testing.T) {
	animal := model.Animal{}

	// 获取 Animal 类型的反射类型
	t := reflect.TypeOf(animal)

	// 创建一个空的 map，用于存储字段名
	fieldMap := make(map[string]bool)

	// 遍历结构体的所有字段
	fmt.Println("All fields:", t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Name

		// 获取 json 标签中的值
		jsonTag := getProcessedJSONTag(field)
		if jsonTag != "" {
			// 如果有 json 标签，则使用标签中的值
			fieldMap[jsonTag] = true
		} else {
			// 如果没有 json 标签，则转换字段名为下划线命名
			convertedFieldName := camelToSnake(fieldName)
			fieldMap[convertedFieldName] = true
		}
	}

	// 打印 map
	fmt.Println("Fields as keys:")
	for key := range fieldMap {
		fmt.Println(key)
	}
}

// 将驼峰命名转换为下划线命名
func camelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}
