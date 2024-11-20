package test

import (
	"fmt"
	"testing"
	"time"
)

// 假设的 EncounterResult 结构体
type EncounterResult struct {
	Title     string
	Content   string
	UpdatedAt time.Time
}

// ToString 方法
func (e EncounterResult) ToString() string {
	return fmt.Sprintf(`路遇笔记标题：%s；路遇笔记内容：%s；最后更新时间：%v`, e.Title, e.Content, e.UpdatedAt)
}

// 测试 EncounterResult 的 ToString 方法
func TestEncounterResult_ToString(t *testing.T) {
	// 设置一个时间点，用于测试
	testTime := time.Now()

	// 创建一个 EncounterResult 实例
	testResult := EncounterResult{
		Title:     "测试笔记",
		Content:   "这是测试笔记的内容",
		UpdatedAt: testTime,
	}

	// 调用 ToString 方法
	resultString := testResult.ToString()

	t.Log("resultString:", resultString)
	// 构建期望的结果字符串
	expectedString := fmt.Sprintf(`路遇笔记标题：%s；路遇笔记内容：%s；最后更新时间：%v`, testResult.Title, testResult.Content, testResult.UpdatedAt)

	// 比较实际结果和期望结果
	if resultString != expectedString {
		t.Errorf("ToString() failed, expected %q, got %q", expectedString, resultString)
	}
}
