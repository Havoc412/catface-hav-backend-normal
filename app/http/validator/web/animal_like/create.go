package animal_like

import (
	"catface/app/global/consts"
	"catface/app/http/controller/web"
	"catface/app/http/validator/core/data_transfer"
	"catface/app/utils/response"

	"github.com/gin-gonic/gin"
)

type Create struct {
	UserId   int `form:"user_id" json:"user_id" binding:"required,min=1"`
	AnimalId int `form:"animal_id" json:"animal_id" binding:"required,min=1"`
}

func (c Create) CheckParams(context *gin.Context) {
	if err := context.ShouldBind(&c); err != nil {
		response.ValidatorError(context, err)
		return
	}
	extraAddBindDataContext := data_transfer.DataAddContext(c, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "AnimalLike Create表单验证器json化失败", "")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&web.AnimalLike{}).Create(extraAddBindDataContext)
	}
}
