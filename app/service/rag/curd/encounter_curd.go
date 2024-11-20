package curd

import (
	"catface/app/model"
	"catface/app/model_es"
	"catface/app/model_res"
)

func CreateEncounterCurdFactory() *EncounterCurd {
	return &EncounterCurd{
		enc:    model.CreateEncounterFactory(""),
		enc_es: model_es.CreateEncounterESFactory(nil),
	}
}

type EncounterCurd struct {
	enc    *model.Encounter
	enc_es *model_es.Encounter
}

func (e *EncounterCurd) TopK(embedding []float64, k int) (temp []model_res.EncounterResult, err error) {
	// ES: TopK
	encounters_es, err := e.enc_es.TopK(embedding, k)
	if err != nil {
		return
	}

	// MySQL 补充信息
	var ids []int64
	for _, encounter := range encounters_es {
		ids = append(ids, encounter.Id)
	}
	encounters := e.enc.ShowByIDs(ids, "id", "title", "content", "updated_at")
	for _, encounter := range encounters {
		for _, encounter_es := range encounters_es {
			if encounter.Id == encounter_es.Id {
				temp = append(temp, *model_res.NewEncounterResult(&encounter, &encounter_es))
			}
		}
	}
	return
}
