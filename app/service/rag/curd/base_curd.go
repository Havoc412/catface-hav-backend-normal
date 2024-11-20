package curd

import (
	"catface/app/global/consts"
	"catface/app/model_res"
	"fmt"
)

func TopK(mode string, embedding []float64, k int) (temp []model_res.DocInterface, err error) {
	switch mode {
	case consts.RagChatModeKnowledge:
		results, err := CreateDocCurdFactory().TopK(embedding, k)
		if err != nil {
			return nil, fmt.Errorf("TopK: 获取知识库TopK失败: %w", err)
		}
		for _, result := range results {
			temp = append(temp, result)
		}

	case consts.RagChatModeDiary:
		results, err := CreateEncounterCurdFactory().TopK(embedding, k)
		if err != nil {
			return nil, fmt.Errorf("TopK: 获取路遇笔记TopK失败: %w", err)
		}
		for _, result := range results {
			temp = append(temp, result)
		}

	default:
		if mode == "" {
			err = fmt.Errorf("TopK: mode不能为空")
		} else {
			err = fmt.Errorf("TopK: 不支持的mode: %s", mode)
		}
	}
	return temp, err
}
