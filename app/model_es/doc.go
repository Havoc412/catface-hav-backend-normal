package model_es

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

// UPDATE 因为 chunck 还是 Python 来处理会比较方便，所以 Go 这边主要还是处理查询相关的操作。
// func (d *Doc) InsertDocument() error {
// 	return nil
// }
