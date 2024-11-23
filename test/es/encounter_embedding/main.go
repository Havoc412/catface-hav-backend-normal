package main

import (
	"fmt"
	"math/rand"
	"strings"
)

func main() {
	embedding := make([]float64, 768)
	for i := range embedding {
		embedding[i] = rand.Float64()
	}

	// 将嵌入向量转换为字符串，每个元素之间用逗号隔开
	embeddingStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(embedding)), ","), "[]")

	// 打印前几个元素以验证
	fmt.Println(embeddingStr)
}
