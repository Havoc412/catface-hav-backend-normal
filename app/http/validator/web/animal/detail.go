package animal

import (
	"catface/app/global/consts"
	"catface/app/http/controller/web"
	"catface/app/http/validator/core/data_transfer"
	"catface/app/utils/response"

	"github.com/gin-gonic/gin"
)

type Detail struct {
	AnmId int64 `form:"anm_id" json:"anm_id"`
	// TODO 暂时没有用到这个模块，GinSK 的架构对于 Path 的处理方式我还不确定，先用临时方案。
}

func (d Detail) CheckParams(context *gin.Context) {
	if err := context.ShouldBind(&d); err != nil {
		response.ValidatorError(context, err)
		return
	}

	extraAddBindDataContext := data_transfer.DataAddContext(d, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "animialList表单验证器json化失败", "")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&web.Animals{}).Detail(extraAddBindDataContext)
	}
}
