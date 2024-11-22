package test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

// EncounterResult2 结构体定义
type EncounterResult2 struct {
	Id        int64      `json:"id" explain:"路遇笔记ID"`
	Title     string     `json:"title" explain:"路遇笔记标题"`
	Content   string     `json:"content" explain:"路遇笔记内容"`
	UpdatedAt *time.Time `json:"updated_at" explain:"最后更新时间"`
}

// StructToString 使用反射将结构体的内容组织为字符串
func StructToString(v interface{}) string {
	val := reflect.ValueOf(v)
	typ := val.Type()

	var result []string
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("explain")
		value := val.Field(i).Interface()
		result = append(result, fmt.Sprintf("%s：%v", tag, value))
	}
	return strings.Join(result, "；")
}

func TestExplain(t *testing.T) {
	// 示例数据
	updatedAt := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)
	encounter := EncounterResult2{
		Id:        1,
		Title:     "遇见小猫",
		Content:   "今天在公园遇到了一只可爱的小猫。",
		UpdatedAt: &updatedAt,
	}

	// 调用 StructToString 函数
	t.Logf("结构体内容：%v", StructToString(encounter))
}
