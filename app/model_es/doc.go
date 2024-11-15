package model_es

import (
	"catface/app/utils/data_bind"
	"catface/app/utils/model_handler"
	"encoding/json"
	"fmt"
)

// INFO @brief 这个文件就是处理 ES 存储文档特征向量的集中处理；

func CreateDocESFactory() *Doc {
	return &Doc{}
}

type Doc struct {
	Id        int64     `json:"id"` // 对应 MySQL 中存储的文档 ID。
	Content   string    `json:"content"`
	Embedding []float64 `json:"embedding"`
}

func (d *Doc) IndexName() string {
	return "catface_docs"
}

func (d *Doc) TopK(embedding []float64, k int) ([]Doc, error) {
	// 将 embedding 数组转换为 JSON 格式
	params := map[string]interface{}{
		"query_vector": embedding,
	}
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	// 构建请求体
	body := fmt.Sprintf(`{
        "size": %d,
        "query": {
            "script_score": {
                "query": {"match_all": {}},
                "script": {
                    "source": "cosineSimilarity(params.query_vector, 'embedding') + 1.0",
                    "params": %s
                }
            }
        },
		"_source": ["content"]
    }`, k, string(paramsJSON))

	hits, err := model_handler.SearchRequest(body, d.IndexName())
	if err != nil {
		return nil, err
	}

	var docs []Doc
	for _, hit := range hits {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"].(map[string]interface{})
		var doc Doc
		if err := data_bind.ShouldBindFormMapToModel(source, &doc); err != nil {
			continue
		}

		docs = append(docs, doc)
	}

	return docs, nil
}

// UPDATE 因为 chunck 还是 Python 来处理会比较方便，所以 Go 这边主要还是处理查询相关的操作。
// func (d *Doc) InsertDocument() error {
// 	return nil
// }
