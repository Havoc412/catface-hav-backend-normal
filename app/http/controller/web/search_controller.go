package web

import (
	"catface/app/global/consts"
	"catface/app/model"
	"catface/app/model_es"
	"catface/app/utils/response"

	"github.com/gin-gonic/gin"
)

type Search struct {
}

/**
 * @description: 全局搜索：AnmName + Encounter
 * @param {*gin.Context} context
 * @return {*}
 */
func (s *Search) SearchAll(context *gin.Context) {
	query := context.GetString(consts.ValidatorPrefix + "query")

	var animals []model.Animal
	var encounters []model.Encounter

	// 1. Animal Name  // TODO 增加字段的过滤，看前端了。
	animals = model.CreateAnimalFactory("").ShowByName(query)

	// 2. Encounter
	_, _ = model_es.CreateEncounterESFactory(nil).QueryDocumentsMatchAll(query)

	// if len(encounterIds) > 0 {
	// 	encounters = model.CreateEncounterFactory("").ShowByIDs(encounterIds)
	// }

	// 3. Knowledge
	knowledges, _ := model_es.CreateKnowledgeESFactory().QueryDocumentsMatchAll(query, 3)

	response.Success(context, consts.CurdStatusOkMsg, gin.H{
		"animals":    animals,
		"encounters": encounters,
		"knowledges": knowledges,
	})
}
