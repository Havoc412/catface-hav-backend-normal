package encounter_like

import (
	"catface/app/global/consts"
	"catface/app/http/controller/web"
	"catface/app/http/validator/core/data_transfer"
	"catface/app/utils/response"

	"github.com/gin-gonic/gin"
)

type Delete struct {
	UserId      int64 `form:"user_id" json:"user_id"`
	EncounterId int64 `form:"encounter_id" json:"encounter_id"`
}

func (d Delete) CheckParams(context *gin.Context) {
	if err := context.ShouldBind(&d); err != nil {
		response.ValidatorError(context, err)
		return
	}

	extraAddBindDataContext := data_transfer.DataAddContext(d, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "EncounterLike Delete 表单验证器json化失败", "")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&web.EncounterLike{}).Delete(extraAddBindDataContext)
	}
}
