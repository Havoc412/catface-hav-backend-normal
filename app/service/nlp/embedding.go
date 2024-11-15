package nlp

import (
	"catface/app/global/variable"
	"catface/app/utils/micro_service"
	"context"

	"github.com/carlmjohnson/requests"
)

type EmbeddingRes struct {
	Status    int       `json:"status"`
	Message   string    `json:"message"`
	Embedding []float64 `json:"embedding"`
}

func GetEmbedding(text string) ([]float64, bool) {
	body := map[string]interface{}{
		"text": text,
	}
	var res EmbeddingRes
	err := requests.URL(micro_service.FetchPythonServiceUrl("rag/bge_embedding")).
		BodyJSON(&body).
		ToJSON(&res).
		Fetch(context.Background())
	if err != nil {
		variable.ZapLog.Error("获取嵌入向量失败: " + err.Error())
	}
	if res.Status != 200 {
		return nil, false
	} else {
		return res.Embedding, true
	}
}
