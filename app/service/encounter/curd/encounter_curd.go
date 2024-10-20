package curd

import "catface/app/model"

func CreateEncounterCurdFactory() *EncounterCurd {
	return &EncounterCurd{model.CreateEncounterFactory("")}
}

type EncounterCurd struct {
	encounter *model.Encounter
}

func (e *EncounterCurd) List(num, skip, user_id int) []model.EncounterList {
	if num == 0 {
		num = 15
	}
	return model.CreateEncounterFactory("").Show(num, skip, user_id)
}
