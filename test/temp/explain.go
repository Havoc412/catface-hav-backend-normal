package main

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// EncounterResult 结构体定义
type EncounterResult struct {
	Id        int64      `json:"id"`
	Title     string     `json:"title" explain:"路遇笔记标题"`
	Content   string     `json:"content" explain:"路遇笔记内容"`
	UpdatedAt *time.Time `json:"updated_at" explain:"最后更新时间"`
	NoTag     string     `json:"no_tag"` // 没有 explain 标签的字段
}

// StructToString 使用反射将结构体的内容组织为字符串，忽略没有 explain 标签的字段
func StructToString(v interface{}) string {
	val := reflect.ValueOf(v)
	typ := val.Type()

	var result []string
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("explain")
		if tag == "" {
			continue // 跳过没有 explain 标签的字段
		}
		value := val.Field(i).Interface()
		result = append(result, fmt.Sprintf("%s：%v", tag, value))
	}
	return strings.Join(result, "；")
}

func main() {
	// 示例数据
	updatedAt := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)
	encounter := EncounterResult{
		Id:        1,
		Title:     "遇见小猫",
		Content:   "今天在公园遇到了一只可爱的小猫。",
		UpdatedAt: &updatedAt,
		NoTag:     "这个字段没有 explain 标签",
	}

	// 调用 StructToString 函数
	fmt.Println(StructToString(encounter))
}
