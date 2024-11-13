package knowledge

import (
	"catface/app/global/consts"
	"catface/app/http/controller/web"
	"catface/app/http/validator/core/data_transfer"
	"catface/app/utils/response"

	"github.com/gin-gonic/gin"
)

type Random struct {
	Num int `form:"num" json:"num"`
}

func (r Random) CheckParams(context *gin.Context) {
	if err := context.ShouldBind(&r); err != nil {
		response.ValidatorError(context, err)
		return
	}

	extraAddBindDataContext := data_transfer.DataAddContext(r, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "Random Knowledge 表单验证器json化失败", "")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&web.Knowledge{}).RandomList(extraAddBindDataContext)
	}
}
