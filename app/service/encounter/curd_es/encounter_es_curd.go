package curd_es

import (
	"catface/app/model"
	"catface/app/model_es"
	"catface/app/service/nlp"
	"catface/app/utils/data_explain"
	"fmt"
)

func CreateEncounterESCurdFactory(encounter *model.Encounter) *EncounterESCurd {
	return &EncounterESCurd{
		model_es.CreateEncounterESFactory(encounter),
	}
}

type EncounterESCurd struct {
	encounter_es *model_es.Encounter
}

func (e *EncounterESCurd) InsertDocument() error {
	var ok bool
	explian := data_explain.GenerateExplainStringForEmbedding(e.encounter_es)
	if e.encounter_es.Embedding, ok = nlp.GetEmbeddingOneString(explian); !ok {
		return fmt.Errorf("nlp embedding service error")
	}

	return e.encounter_es.InsertDocument()
}
