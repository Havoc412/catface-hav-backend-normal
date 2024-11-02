package model

import "gorm.io/gorm"

func CreateEncounterAnimalLinkFactory(sqlType string) *EncounterAnimalLink {
	return &EncounterAnimalLink{DB: UseDbConn(sqlType)}
}

/**
 * @description: 路遇 & 关联毛茸茸 的关联表，将 []int 拆为 int 的子表，方便查询。
 * @return {*}
 */
type EncounterAnimalLink struct {
	*gorm.DB    `gorm:"-" json:"-"`
	EncounterId int64 `gorm:"column:encounter_id;index" json:"encounter_id"`
	Encounter   Encounter
	AnimalId    int64 `gorm:"column:animal_id;index" json:"animal_id"`
	Animal      Animal
}

func (e *EncounterAnimalLink) Insert(encounterId int64, animalId []float64) bool {
	// Build Slice
	var results []EncounterAnimalLink
	for _, id := range animalId {
		results = append(results, EncounterAnimalLink{EncounterId: encounterId, AnimalId: int64(id)})
	}

	// 调用批处理插入方法
	return e.Create(&results).Error == nil
}

func (e *EncounterAnimalLink) ShowByEncounterId(encounterId int64) ([]int64, bool) {
	var results []EncounterAnimalLink
	if err := e.Where("encounter_id = ?", encounterId).Find(&results).Error; err != nil {
		// 处理错误情况，例如日志记录或返回错误信息
		return nil, false
	}
	intSlice := make([]int64, len(results))
	for i, result := range results {
		intSlice[i] = result.AnimalId // 假设 AnimalId 是你需要提取的字段
	}
	return intSlice, true
}

func (e *EncounterAnimalLink) ShowByEncounterIdFirst(encounterId int64) (int64, bool) {
	var results EncounterAnimalLink
	if err := e.Where("encounter_id = ?", encounterId).First(&results).Error; err != nil {
		// 处理错误情况，例如日志记录或返回错误信息
		return 0, false
	}
	return results.AnimalId, true
}
