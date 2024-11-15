package model_es

// INFO @brief 这个文件就是处理 ES 存储文档特征向量的集中处理

func CreateDocESFactory() *Doc {
	return &Doc{}
}

type Doc struct {
	Content   string    `json:"content"`
	Embedding []float64 `json:"embedding"`
}

func (d *Doc) IndexName() string {
	return "catface_docs"
}

func (d *Doc) InsertDocument() error {
	return nil
}
