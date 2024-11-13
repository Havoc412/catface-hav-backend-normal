package curd

import (
	"catface/app/model"
	"catface/app/model_es"
	"catface/app/utils/query_handler"
	"strconv"
)

func CreateEncounterCurdFactory() *EncounterCurd {
	return &EncounterCurd{model.CreateEncounterFactory("")}
}

type EncounterCurd struct {
	encounter *model.Encounter
}

func (e *EncounterCurd) List(num, skip, user_id int, mode string) (result []model.EncounterList) {
	if num == 0 {
		num = 10
	}

	var likedAnimalIds []int
	switch mode {
	case "liked":
		likedAnimalIds = model.CreateAnimalLikeFactory("").LikedCats(user_id)
	}
	result = model.CreateEncounterFactory("").Show(num, skip, user_id, likedAnimalIds)
	return
}

func (e *EncounterCurd) Detail(id string) *model.EncounterDetail {
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

	// 1.1 处理 Photos 为 []string，同时忽略原本的 Photos 字段。
	encounter.PhotosList = query_handler.StringToStringArray(encounter.Photos)
	encounter.Photos = "" // 清空。

	// 2. user data
	user, err := model.CreateUserFactory("").ShowByID(encounter.UsersModelId, "user_avatar", "user_name", "id")
	if err != nil {
		return nil
	}

	// 3. animals data
	var animals []model.Animal
	if animals_id, ok := model.CreateEncounterAnimalLinkFactory("").ShowByEncounterId(encounter.Id); ok {
		animals = model.CreateAnimalFactory("").ShowByIDs(animals_id, "avatar", "name", "id")
	}

	// 4. 合并
	return &model.EncounterDetail{
		Encounter:  *encounter,
		UsersModel: *user,
		Animals:    animals,
	}
}

func (e *EncounterCurd) MatchAll(query string, num int) (tmp []model.Encounter) {
	// 1. encounter ES
	encountersFromES, err := model_es.CreateEncounterESFactory(nil).QueryDocumentsMatchAll(query, num)
	if err != nil || len(encountersFromES) == 0 {
		return nil
	}

	var ids []int64
	for _, encounter := range encountersFromES {
		ids = append(ids, encounter.Id)
	}

	// 2. encounter SQL
	encountersFromSQL := model.CreateEncounterFactory("").ShowByIDs(ids, "id", "avatar", "user_id")

	// 3. users
	ids = nil
	for _, encounter := range encountersFromSQL {
		ids = append(ids, encounter.UsersModelId)
	}
	users := model.CreateUserFactory("").ShowByIDs(ids, "user_avatar", "user_name", "id")

	// end. Merge
	for _, enencountersFromES := range encountersFromES {
		for _, encounter := range encountersFromSQL {
			for _, user := range users {
				if encounter.Id == enencountersFromES.Id && encounter.UsersModelId == user.Id {
					encounter.TagsList = enencountersFromES.Tags
					encounter.TagsHighlight = enencountersFromES.TagsHighlight
					encounter.Title = enencountersFromES.Title
					encounter.Content = enencountersFromES.Content

					encounter.UsersModel = &user
					tmp = append(tmp, encounter)
				}
			}
		}
	}

	return
}
