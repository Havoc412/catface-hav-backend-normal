package curd

import (
	"catface/app/model"

	"github.com/gin-gonic/gin"
)

func CreateDocCurdFactory() *DocCurd {
	return &DocCurd{doc: model.CreateDocFactory("")}
}

type DocCurd struct {
	doc *model.Doc
}

// UPDATE 好像有点没有必要性。
func (d *DocCurd) InsertDocumentData(c *gin.Context) bool {
	// TODO insert embedding to ES // INFO 调用 py 的服务？
	// STAGE insert data to mysql
	return d.doc.InsertDocumentData(c)
}
