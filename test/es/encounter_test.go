package test

import (
	"catface/app/global/variable"
	"catface/app/model"
	"catface/app/model_es"
	_ "catface/bootstrap"
	"fmt"
	"testing"
)

func TestEncounterEs(t *testing.T) {
	// 示例数据
	encounterOri := &model.Encounter{
		BaseModel: model.BaseModel{
			Id: 4,
		},
		Title:     "猪皮伤势轻，需静养猪皮伤势轻，需静养",
		Content:   "猪皮被带到医院检查了，拍片结果显示损伤不严重，静养即可自愈。建议这段时间不要折腾他，让老登好好休息。",
		TagsSlice: []string{"猪皮", "脚伤", "骗保"},
	}

	encounter := model_es.CreateEncounterESFactory(encounterOri)

	// 插入文档
	// if err := encounter.InsertDocument(); err != nil {
	// 	t.Fatalf("插入文档时出错: %s", err)
	// }
	go encounter.InsertDocument()

	// // 更新文档
	// encounter.Content = "更新: 猪皮被带到医院检查了，拍片结果显示损伤不严重，静养即可自愈。建议这段时间不要折腾他，让老登好好休息。"
	// if err := encounter.UpdateDocument(variable.ElasticClient, encounter); err != nil {
	// 	t.Fatalf("更新文档时出错: %s", err)
	// }

	// 查询文档
	encounters, err := encounter.QueryDocumentsMatchAll(variable.ElasticClient, "猪皮")
	if err != nil {
		t.Fatalf("查询文档时出错: %s", err)
	}

	// for _, e := range encounters {
	// 	fmt.Printf("ID: %d, 标题: %s, 内容: %s, 标签: %v\n", e.Id, e.Title, e.Content, e.Tags)
	// }

	for _, e := range encounters {
		fmt.Printf("ID: %d\n", e)
	}
}
