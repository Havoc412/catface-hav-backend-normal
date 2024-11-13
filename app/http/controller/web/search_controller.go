package web

import (
	"catface/app/global/consts"
	"catface/app/model_es"
	animal_curd "catface/app/service/animals/curd"
	encouner_curd "catface/app/service/encounter/curd"
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

	// 1. Animal Name  // TODO 增加字段的过滤，看前端了。
	// animals = model.CreateAnimalFactory("").ShowByName(query)
	animals := animal_curd.CreateAnimalsCurdFactory().MatchAll(query, 3)

	// 2. Encounter
	encounters := encouner_curd.CreateEncounterCurdFactory().MatchAll(query, 3)

	// 3. Knowledge
	knowledges, _ := model_es.CreateKnowledgeESFactory().QueryDocumentsMatchAll(query, 3)

	response.Success(context, consts.CurdStatusOkMsg, gin.H{
		"animals":    animals,
		"encounters": encounters,
		"knowledges": knowledges,
	})
}
