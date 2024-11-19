package model_res

import (
	"catface/app/model"
	"catface/app/model_es"
	"time"
)

// BUG 存在 依賴循環
func NewDocResult(doc *model.Doc, doc_es *model_es.Doc) *DocResult {
	return &DocResult{
		Type:      "doc",
		Id:        doc.Id,
		Name:      doc.Name,
		Content:   doc_es.Content,
		UpdatedAt: doc.UpdatedAt,
	}
}

type DocResult struct {
	Type      string     `json:"type"`
	Id        int64      `json:"id"`
	Name      string     `json:"name"`
	Content   string     `json:"content"`
	UpdatedAt *time.Time `json:"updated_at"`
}
