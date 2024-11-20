package data_explain

import (
	"fmt"
	"reflect"
	"strings"
)

/**
 * @description: 集成 Struct -> Explain for RAG；
 * @param {interface{}} v
 * @return {*}
 */
func GenerateExplainStringForEmbedding(v interface{}) string {
	val := reflect.ValueOf(v)
	typ := val.Type()

	var result []string
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("explain")
		if tag == "" {
			continue
		}
		value := val.Field(i).Interface()
		result = append(result, fmt.Sprintf("%s：%v", tag, value))
	}
	return strings.Join(result, "；")
}
