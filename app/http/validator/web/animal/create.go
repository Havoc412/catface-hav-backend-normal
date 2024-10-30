package animal

import (
	"catface/app/global/consts"
	"catface/app/http/controller/web"
	"catface/app/http/validator/common/location"
	"catface/app/http/validator/core/data_transfer"
	"catface/app/utils/response"

	"github.com/gin-gonic/gin"
)

type Extra struct {
	NickNames []string `form:"nick_names" json:"nick_names"`
	Tags      []string `form:"tags" json:"tags"`
}

type Health struct {
	Sterilization uint8 `form:"sterilization" json:"sterilization"`
	Vaccination   uint8 `form:"vaccination" json:"vaccination"`
	Deworming     uint8 `form:"deworming" json:"deworming"`
}

type Create struct {
	Name        string   `form:"name" json:"name" binding:"required"`
	Breed       uint8    `form:"breed" json:"breed" binding:"required"`
	Gender      uint8    `form:"gender" json:"gender" binding:"required"`
	Status      uint8    `form:"status" json:"status" binding:"required"`
	Description string   `form:"description" json:"description"`
	Birthday    string   `form:"birthday" json:"birthday"`
	Photos      []string `form:"photos" json:"photos"`

	Health Health       `form:"health" json:"health"`
	UserId int          `form:"user_id" json:"user_id" binding:"required,numeric"`
	Poi    location.Poi `form:"poi" json:"poi"`
	Extra  Extra        `form:"extra" json:"extra"`

	FaceModelScore float64 `form:"face_model_score" json:"face_model_score"`

	Mode bool `form:"mode" json:"mode"` // INFO 0 default；1 缓存的效果。
}

func (c Create) CheckParams(context *gin.Context) {
	if err := context.ShouldBind(&c); err != nil {
		response.ValidatorError(context, err)
		return
	}

	extraAddBindDataContext := data_transfer.DataAddContext(c, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "Animal Create 表单验证器json化失败", "")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&web.Animals{}).Create(extraAddBindDataContext)
	}
}
