package curd

import (
	"catface/app/model"
	"catface/app/model_es"
	"catface/app/model_res"
)

func CreateDocCurdFactory() *DocCurd {
	return &DocCurd{
		doc:    model.CreateDocFactory(""),
		doc_es: model_es.CreateDocESFactory()}
}

type DocCurd struct { // 组合方法的使用
	doc    *model.Doc
	doc_es *model_es.Doc
}

func (d *DocCurd) TopK(embedding []float64, k int) (temp []model_res.DocResult, err error) {
	// ES：TopK
	docs_es, err := d.doc_es.TopK(embedding, k)
	if err != nil {
		return
	}

	// MySQL：补充基本信息
	var ids []int64
	for _, doc := range docs_es {
		ids = append(ids, doc.Id)
	}
	docs := d.doc.ShowByIds(ids, "id", "name")

	// 装载
	for _, doc := range docs {
		for _, doc_es := range docs_es {
			if doc.Id == doc_es.Id {
				temp = append(temp, *model_res.NewDocResult(&doc, &doc_es))
			}
		}
	}

	return
}
