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
	EncounterId int `gorm:"column:encounter_id;index" json:"encounter_id"`
	Encounter   Encounter
	AnimalId    int `gorm:"column:animal_id;index" json:"animal_id"`
	Animal      Animal
}

func (e *EncounterAnimalLink) Insert(encounterId int, animalId []float64) bool {
	// Build Slice
	var results []EncounterAnimalLink
	for _, id := range animalId {
		results = append(results, EncounterAnimalLink{EncounterId: encounterId, AnimalId: int(id)})
	}

	// 调用批处理插入方法
	return e.Create(&results).Error == nil
}
