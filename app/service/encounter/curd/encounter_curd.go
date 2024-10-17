package curd

import "catface/app/model"

func CreateEncounterCurdFactory() *EncounterCurd {
	return &EncounterCurd{model.CreateEncounterFactory("")}
}

type EncounterCurd struct {
	encounter *model.Encounter
}

func (e *EncounterCurd) Store(encounter *model.Encounter) bool {

	return model.CreateEncounterFactory()
}
