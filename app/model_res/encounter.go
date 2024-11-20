package model_res

import (
	"catface/app/model"
	"catface/app/model_es"
	"fmt"
	"time"
)

func NewEncounterResult(encounter *model.Encounter, encounter_es *model_es.Encounter) *EncounterResult {
	return &EncounterResult{
		DocBase:   DocBase{Type: "encounter"},
		Id:        encounter.Id,
		Title:     encounter.Title,
		Content:   encounter.Content,
		UpdatedAt: encounter.UpdatedAt}
}

type EncounterResult struct {
	DocBase
	Id        int64      `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (e EncounterResult) ToString() string {
	return fmt.Sprintf(`路遇笔记标题：%s；路遇笔记内容：%s；`, e.Title, e.Content)
}
