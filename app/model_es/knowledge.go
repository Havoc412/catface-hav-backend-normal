package model_es

import (
	"bytes"
	"catface/app/global/variable"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func CreateKnowledgeESFactory() *Knowledge {
	return &Knowledge{}
}

type Knowledge struct {
	Id      int32    `json:"id"` // TIP int64 会炸 ES 的 ‘integer’
	Dirs    []string `json:"dirs"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
}

func (k *Knowledge) IndexName() string {
	return "catface_knowledges"
}

func (k *Knowledge) InsertDocument() error {
	ctx := context.Background()

	k.Id = int32(time.Now().UnixNano() / 1e6) // 将纳秒级时间戳转换为毫秒级 // INFO 自动补充时间戳为 ID

	// 将结构体转换为 JSON 字符串
	data, err := json.Marshal(k)
	if err != nil {
		return err
	}

	// 创建请求
	req := esapi.IndexRequest{
		Index: k.IndexName(),
		// DocumentID: fmt.Sprintf("%d", k.Id),
		Body:    bytes.NewReader(data),
		Refresh: "true",
	}

	// 发送请求
	res, err := req.Do(ctx, variable.ElasticClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var k map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&k); err != nil {
			return fmt.Errorf("error parsing the response body: %s", err)
		} else {
			return fmt.Errorf("[%s] %s: %s",
				res.Status(),
				k["error"].(map[string]interface{})["type"],
				k["error"].(map[string]interface{})["reason"],
			)
		}
	}

	return nil
}

// RandomDocuments 随机查询 num 个文档
func (k *Knowledge) RandomDocuments(num int) ([]*Knowledge, error) {
	ctx := context.Background()

	// 创建本地随机数生成器  // TIP rand.Seed() 在 Go1.20 之后弃用了。
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 创建请求
	req := esapi.SearchRequest{
		Index: []string{k.IndexName()},
		Body: strings.NewReader(fmt.Sprintf(`{
			"size": %d,
			"query": {
				"function_score": {
					"query": { "match_all": {} },
					"random_score": {
						"seed": %d
					}
				}
			}
		}`, num, rng.Int63())),
	}

	// 发送请求
	res, err := req.Do(ctx, variable.ElasticClient)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var k map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&k); err != nil {
			return nil, fmt.Errorf("error parsing the response body: %s", err)
		} else {
			return nil, fmt.Errorf("[%s] %s: %s",
				res.Status(),
				k["error"].(map[string]interface{})["type"],
				k["error"].(map[string]interface{})["reason"],
			)
		}
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	// 提取文档
	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})

	documents := make([]*Knowledge, len(hits))
	for i, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})

		// 将 []interface{} 转换为 []string
		dirs := make([]string, len(source["dirs"].([]interface{})))
		for j, dir := range source["dirs"].([]interface{}) {
			dirs[j] = dir.(string)
		}

		doc := &Knowledge{
			Dirs:    dirs,
			Title:   source["title"].(string),
			Content: source["content"].(string),
		}
		documents[i] = doc
	}

	return documents, nil
}
