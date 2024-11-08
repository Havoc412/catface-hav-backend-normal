package animal

import (
	"catface/app/global/consts"
	"catface/app/http/controller/web"
	"catface/app/http/validator/core/data_transfer"
	"catface/app/utils/response"

	"github.com/gin-gonic/gin"
)

type List struct {
	Attrs         string `form:"attrs" json:"attrs"`
	Gender        string `form:"gender" json:"gender"`
	Breed         string `form:"breed" json:"breed"`
	Sterilization string `form:"sterilization" json:"sterilization"`
	Status        string `form:"status" json:"status"`
	Department    string `form:"department" json:"department"`
	Num           int    `form:"num" json:"num"`
	Skip          int    `form:"skip" json:"skip"`
	UserId        int    `form:"user_id" json:"user_id"`

	Mode string `form:"mode" json:"mode"` // INFO 控制 animal_ctl 的查询模式。 // default: 简单调用 List 函数 || prefer: 优先返回和用户关联度更高的目标。
	Key  int64  `form:"key" json:"key"`   // redis の key 值。
}

func (l List) CheckParams(context *gin.Context) {
	if err := context.ShouldBind(&l); err != nil {
		// 将表单参数验证器出现的错误直接交给错误翻译器统一处理即可
		response.ValidatorError(context, err)
		return
	}

	// 该函数主要是将本结构体的字段（成员）按照 consts.ValidatorPrefix+ json标签对应的 键 => 值 形式绑定在上下文，便于下一步（控制器）可以直接通过 context.Get(键) 获取相关值
	extraAddBindDataContext := data_transfer.DataAddContext(l, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "animialList表单验证器json化失败", "")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&web.Animals{}).List(extraAddBindDataContext)
	}
}
