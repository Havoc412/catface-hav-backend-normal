package curd

import (
	"catface/app/model"
	"catface/app/utils/query_handler"
	"strconv"
)

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

func (e *EncounterCurd) Detail(id string) *model.Encounter {
	// 0. check id
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil
	}

	// 1. encounter data
	encounter, err := model.CreateEncounterFactory("").ShowByID(int64(idInt))
	if err != nil {
		return nil
	}

	// 2. user data
	user, err := model.CreateUserFactory("").ShowByID(encounter.UsersModelId, "user_avatar", "user_name")
	if err != nil {
		return nil
	}
	_ = user

	// 3. animals data
	animals_id := query_handler.StringToint64Array(encounter.AnimalsId)
	animals := model.CreateAnimalFactory("").ShowByIDs(animals_id, "avatar", "name")
	_ = animals

	// 4. 合并

	return nil
}
