package test

import (
	"catface/app/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnmFaceBreed(t *testing.T) {
	Init()

	err := DB.AutoMigrate(&model.AnmFaceBreed{})
	if err != nil {
		t.Error(err)
	}

	// // INFO 查询表上的所有索引
	// var indexes []struct {
	// 	IndexName  string
	// 	ColumnName string
	// }
	// DB.Raw(`SELECT INDEX_NAME, COLUMN_NAME FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?`, "Hav'sCats", "anm_face_breeds").Scan(&indexes)
	// fmt.Println("All Indexes:", len(indexes)) // QUESTION 输出 0 ?
	// for _, index := range indexes {
	// 	fmt.Printf("Index Name: %s, Column Name: %s\n", index.IndexName, index.ColumnName)
	// }

	animalFaceBreed := model.AnmFaceBreed{
		AnimalId: 1,
		Top1:     3,
		Prob1:    0.9,
		Top2:     4,
		Prob2:    0.05,
		Top3:     5,
		Prob3:    0.05,
	}

	// res := DB.Create(&animalFaceBreed)
	// assert.Nil(t, res.Error)

	// 可以进一步检查数据是否正确插入,例如通过查询数据库来验证
	var temp model.AnmFaceBreed
	result := DB.First(&temp, 1) //animalFaceBreed.BaseModel.Id)  // ATT 这里用 Id 直接去拿到默认值 0
	assert.Nil(t, result.Error)
	assert.Equal(t, animalFaceBreed.Top1, temp.Top1)
}
