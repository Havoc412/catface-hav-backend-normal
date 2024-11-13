package model_es

import (
	"bytes"
	"catface/app/global/variable"
	"catface/app/model"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func CreateEncounterESFactory(encounter *model.Encounter) *Encounter {
	if encounter == nil { // UPDATE 这样写好丑。
		return &Encounter{}
	}

	// 我把数值绑定到了工厂创建当中。
	return &Encounter{
		Id:      encounter.Id,
		Title:   encounter.Title,
		Content: encounter.Content,
		Tags:    encounter.TagsList, // TODO 暂时没有对此字段的查询。
	}
}

// INFO 存储能够作为索引存在的数据。
type Encounter struct {
	Id      int64    `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (e *Encounter) IndexName() string {
	return "catface_encounters"
}

func (e *Encounter) InsertDocument() error {
	ctx := context.Background()

	// 将结构体转换为 JSON 字符串
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	// 创建请求
	req := esapi.IndexRequest{
		Index:      e.IndexName(),
		DocumentID: fmt.Sprintf("%d", e.Id),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	// 发送请求
	res, err := req.Do(ctx, variable.ElasticClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return fmt.Errorf("error parsing the response body: %s", err)
		} else {
			return fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	return nil
}

// TODO 改正，仿 Insert
func (e *Encounter) UpdateDocument(client *elasticsearch.Client, encounter *Encounter) error {
	ctx := context.Background()

	// 将结构体转换为 JSON 字符串
	data, err := json.Marshal(map[string]interface{}{
		"doc": encounter,
	})
	if err != nil {
		return err
	}

	// 创建请求
	req := esapi.UpdateRequest{
		Index:      encounter.IndexName(),
		DocumentID: fmt.Sprintf("%d", encounter.Id),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	// 发送请求
	res, err := req.Do(ctx, client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return fmt.Errorf("error parsing the response body: %s", err)
		} else {
			return fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	return nil
}

/**
 * @description: 粗略地包含各种关键词匹配，
 * @param {*elasticsearch.Client} client
 * @param {string} query
 * @return {*} 对应 Encounter 的 id，然后交给 MySQL 来查询详细的信息？
 */
func (e *Encounter) QueryDocumentsMatchAll(query string) ([]int64, error) {
	ctx := context.Background()

	// 创建查询请求
	req := esapi.SearchRequest{ // UPDATE 同时实现查询高亮？
		Index: []string{e.IndexName()},
		Body: strings.NewReader(fmt.Sprintf(`{
			"_source": ["id"],
			"query": {
				"bool": {
					"should": [
						{
							"match": {
								"title": "%s"
							}
						},
						{
							"match": {
								"content": "%s"
							}
						}
					]
				}
			}
		}`, query, query)),
	}

	// 发送请求
	res, err := req.Do(ctx, variable.ElasticClient)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("error parsing the response body: %s", err)
		} else {
			return nil, fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	// 解析响应
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	// 提取命中结果
	hits, ok := r["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("error extracting hits from response")
	}

	fmt.Println(hits)

	// 转换为 id 切片
	var ids []int64
	for _, hit := range hits {
		hitMap := hit.(map[string]interface{})["_source"].(map[string]interface{})
		id := int64(hitMap["id"].(float64))
		ids = append(ids, id)
	}

	return ids, nil

	// // 转换为 Encounter 切片
	// var encounters []*Encounter
	// for _, hit := range hits {
	// 	hitMap := hit.(map[string]interface{})
	// 	source := hitMap["_source"].(map[string]interface{})

	// 	// TIP 将 []interface{} 转换为 []string
	// 	tagsInterface := source["tags"].([]interface{})
	// 	tags := make([]string, len(tagsInterface))
	// 	for i, tag := range tagsInterface {
	// 		tags[i] = tag.(string)
	// 	}

	// 	encounter := &Encounter{
	// 		Id:      int64(source["id"].(float64)),
	// 		Title:   source["title"].(string),
	// 		Content: source["content"].(string),
	// 		Tags:    tags,
	// 	}
	// 	encounters = append(encounters, encounter)
	// }

	// return encounters, nil
}
