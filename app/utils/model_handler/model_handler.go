package model_handler

import (
	"reflect"
	"strings"
	"unicode"
)

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

// 获取处理后的 json 标签值
func getProcessedJSONTag(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag != "" {
		// INFO 移除 omitempty
		jsonTag = strings.Replace(jsonTag, ",omitempty", "", -1)
		jsonTag = strings.Replace(jsonTag, "omitempty", "", -1)
	}
	return jsonTag
}

func GetModelField(v interface{}) map[string]bool {
	t := reflect.TypeOf(v) // TODO 特化处理掉 BaseModel 这样的继承字段

	fieldMap := make(map[string]bool)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Name

		// 获取 json 标签中的值
		jsonTag := getProcessedJSONTag(field)
		if jsonTag != "" {
			// 如果有 json 标签，则使用标签中的值
			fieldMap[jsonTag] = false
		} else {
			// 如果没有 json 标签，则转换字段名为下划线命名
			convertedFieldName := camelToSnake(fieldName)
			fieldMap[convertedFieldName] = false
		}
	}
	fieldMap["id"] = true // INFO default ID 默认都会返回。
	return fieldMap
}
