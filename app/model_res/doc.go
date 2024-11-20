package model_res

import (
	"catface/app/model"
	"catface/app/model_es"
	"time"
)

// INFO 由于直接放到 model 中会导致循环引用，所以放到 model_res 中
func NewDocResult(doc *model.Doc, doc_es *model_es.Doc) *DocResult {
	return &DocResult{
		DocBase:   DocBase{Type: "doc"},
		Id:        doc.Id,
		Name:      doc.Name,
		Content:   doc_es.Content,
		UpdatedAt: doc.UpdatedAt,
	}
}

type DocResult struct {
	DocBase
	Id        int64      `json:"id"`
	Name      string     `json:"name" explain:"文档名称"`
	Content   string     `json:"content" explain:"文档内容"`
	UpdatedAt *time.Time `json:"updated_at" explain:"最后更新时间"`
}

// GetType implements DocInterface.
func (d DocResult) GetType() string {
	panic("unimplemented")
}

/**
 * @description: 实现 DocInterface 接口，输出作为 LLM 的参考内容。
 * @return {*}
 */
func (d DocResult) ToString() string {
	return d.Content
}
