package curd

import (
	"catface/app/global/consts"
	"catface/app/utils/data_explain"
	"fmt"
	"strconv"
)

/**
 * @description: 作为搜索到的文档的集合，目前都是单一的文档类型；// TODO 如何更好的支持多文档的 TopK ？
 * 相当于 DocumentHub 的构造函数。
 * @param {string} mode
 * @param {[]float64} embedding
 * @param {int} k
 * @return {*}
 */
func TopK(mode string, embedding []float64, k int) (dochub DocumentHub, err error) {
	switch mode {
	case consts.RagChatModeKnowledge:
		results, errTemp := CreateDocCurdFactory().TopK(embedding, k)
		if errTemp != nil {
			err = fmt.Errorf("TopK: 获取知识库TopK失败: %w", errTemp)
		}
		for _, result := range results {
			dochub.Docs = append(dochub.Docs, result)
		}

	case consts.RagChatModeDiary:
		results, errTemp := CreateEncounterCurdFactory().TopK(embedding, k)
		if errTemp != nil {
			err = fmt.Errorf("TopK: 获取路遇笔记TopK失败: %w", errTemp)
		}
		for _, result := range results {
			dochub.Docs = append(dochub.Docs, result)
		}

	default:
		if mode == "" {
			err = fmt.Errorf("TopK: mode不能为空")
		} else {
			err = fmt.Errorf("TopK: 不支持的mode: %s", mode)
		}
	}
	return
}

type DocumentHub struct {
	Docs []interface{}
}

func (d *DocumentHub) Length() int {
	return len(d.Docs)
}

func (d *DocumentHub) Explain4LLM() (explain string) {
	for id, doc := range d.Docs {
		explain += strconv.Itoa(id) + "." + data_explain.GenerateExplainStringForEmbedding(doc) + "\n"
	}
	return
}
