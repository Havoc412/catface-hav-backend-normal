package web

import (
	"catface/app/global/consts"
	"catface/app/service/encounter/curd"
	"catface/app/utils/response"

	"github.com/gin-gonic/gin"
)

type Encounters struct {
}

func (e *Encounters) Store(context *gin.Context) {
	userId := context.GetFloat64(consts.ValidatorPrefix + "user_id")
	animalsID := context.GetString(consts.ValidatorPrefix + "animals_id")
	title := context.GetString(consts.ValidatorPrefix + "title")
	context := context.GetString(consts.ValidatorPrefix + "content")
	photos := context.GetString(consts.ValidatorPrefix + "photos")
	laitude := context.GetFloat64(consts.ValidatorPrefix + "latitude")
	longitude := context.GetFloat64(consts.ValidatorPrefix + "longitude")

	encounters := curd.CreateEncounterCurdFactory().Create()
	if encounters == nil {
		response.Success(context, consts.CurdStatusOkMsg, encounters)
	} else {
		response.Fail(context, consts.CurdCreatFailCode, consts.CurdCreatFailMsg, "")
	}
}
