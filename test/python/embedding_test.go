package nlp

import (
	"catface/app/service/nlp"
	_ "catface/bootstrap"
	"fmt"
	"testing"
)

func TestEmbeddingApi(t *testing.T) {
	res, ok := nlp.GetEmbedding([]string{"一段测试文本。"})
	if !ok {
		t.Error("获取嵌入向量失败")
	}
	fmt.Println(res)
}
