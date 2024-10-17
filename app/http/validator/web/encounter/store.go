package encounter

import (
	"catface/app/global/consts"
	"catface/app/http/controller/web"
	"catface/app/http/validator/core/data_transfer"
	"catface/app/utils/response"

	"github.com/gin-gonic/gin"
)

type Store struct {
	UserId    int    `form:"user_id" json:"user_id" binding:"required,numeric"`
	AnimalsId string `form:"animals_id" json:"animals_id" binding:"required"`
	Title     string `form:"title" json:"title" binding:"required"`
	Content   string `form:"content" json:"content"`

	// Avatar string `form:"avatar" json:"avatar"`
	Photos    string  `form:"photos" json:"photos"` // INFO 如果 Photo 为空，那就选取 Animals 的 Avatar
	Laitude   float64 `form:"latitude" json:"latitude"`
	Longitude float64 `form:"longitude" json:"longitude"`
}

func (e Store) CheckParams(context *gin.Context) {
	if err := context.ShouldBind(&e); err != nil {
		response.ValidatorError(context, err)
		return
	}

	extraAddBindDataContext := data_transfer.DataAddContext(e, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "EncounterStore表单验证器json化失败", "")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&web.Encounters{}).Store(extraAddBindDataContext)
	}
}
