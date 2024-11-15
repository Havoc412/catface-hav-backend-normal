package model

import (
	"catface/app/global/variable"
	"catface/app/utils/data_bind"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// INFO @brief 这个 model 是便于宏观管理知识库文件的。

func CreateDocFactory(sqlType string) *Doc {
	return &Doc{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type Doc struct {
	BaseModel
	Name string `gorm:"name" json:"name"` // 文件名保存原本的设定，但是实际存储的【真名】还是借助 Snow + MD5 防止冲突；
	Path string `gorm:"path" json:"path"`
}

func (d *Doc) TableName() string { return "docs" }

func (d *Doc) InsertDocumentData(c *gin.Context) (int64, bool) {
	var tmp Doc
	if err := data_bind.ShouldBindFormDataToModel(c, &tmp); err == nil {
		if res := d.Create(&tmp); res.Error == nil {
			return tmp.Id, true
		} else {
			variable.ZapLog.Error("Doc 数据新增出错", zap.Error(res.Error))
		}
	} else {
		variable.ZapLog.Error("Doc 数据绑定出错", zap.Error(err))
	}
	return 0, false
}
