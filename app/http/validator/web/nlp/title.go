package nlp

import (
	"catface/app/global/consts"
	"catface/app/http/controller/web"
	"catface/app/http/validator/core/data_transfer"
	"catface/app/utils/response"

	"github.com/gin-gonic/gin"
)

// 用于生成标题的文本素材
type Title struct {
	Content     string   `form:"content" json:"content"`
	Title       string   `form:"title" json:"title"` // 原标题
	Tags        []string `form:"tags" json:"tags"`
	AnimalsName []string `form:"animals_name" json:"animals_name"`
}

func (t Title) CheckParams(context *gin.Context) {
	if err := context.ShouldBind(&t); err != nil {
		response.ValidatorError(context, err)
		return
	}

	extraAddBindDataContext := data_transfer.DataAddContext(t, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "Animal Create 表单验证器json化失败", "")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&web.Nlp{}).Title(extraAddBindDataContext)
	}

}
