package model_es

import (
	"bytes"
	"catface/app/global/consts"
	"catface/app/global/variable"
	"catface/app/model"
	"catface/app/utils/data_bind"
	"catface/app/utils/model_handler"
	"context"
	"encoding/json"
	"fmt"

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

	TagsHighlight []string `json:"tags_highlight"`
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
func (e *Encounter) QueryDocumentsMatchAll(query string, num int) ([]Encounter, error) {
	body := fmt.Sprintf(`{
  "size": %d, 
  "query": {
    "bool": {
      "should": [
        {"match": {"tags": "%s"}},
        {"match": {"content": "%s"}},
        {"match": {"title": "%s"}}
      ]
    }
  },
  "highlight": {
    "pre_tags": ["%v"],
    "post_tags": ["%v"],
    "fields": {
      "title": {},
      "content": {
        "fragment_size" : 15
      },
      "tags": {
        "pre_tags": [""],
        "post_tags": [""]
      }
    }
  }
}`, num, query, query, query, consts.PreTags, consts.PostTags)

	hits, err := model_handler.SearchRequest(body, e.IndexName())
	if err != nil {
		return nil, err
	}

	var encounters []Encounter
	for _, hit := range hits {
		data := model_handler.MergeSouceWithHighlight(hit.(map[string]interface{}))

		var encounter Encounter
		if err := data_bind.ShouldBindFormMapToModel(data, &encounter); err != nil {
			continue
		}

		encounters = append(encounters, encounter)
	}

	return encounters, nil
}
