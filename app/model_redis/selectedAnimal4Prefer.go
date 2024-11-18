package model_redis

import (
	"catface/app/global/variable"
	"catface/app/utils/redis_factory"
	"encoding/json"
)

// INFO 辅助 animal - list - prefer 模式下的查询特化。

type SelectedAnimal4Prefer struct {
	EncounteredCatsId []int64 `json:"encountered_cats_id"` // #1 对应第一阶段：近期路遇关联
	NewCatsId         []int64 `json:"new_cats_id"`         // #2 对应第二阶段：近期新增

	key           int64 `json:"-"` // redis 的 key 值
	notPassNewCat bool  `json:"-"`
}

func (s *SelectedAnimal4Prefer) NotPassNew() bool {
	return s.notPassNewCat
}

func (s *SelectedAnimal4Prefer) Length() int {
	return len(s.NewCatsId) + len(s.EncounteredCatsId)
}

func (s *SelectedAnimal4Prefer) NumEnc() int {
	return len(s.EncounteredCatsId)
}

func (s *SelectedAnimal4Prefer) GetAllIds() []int64 {
	return append(s.EncounteredCatsId, s.NewCatsId...)
}

func (s *SelectedAnimal4Prefer) AppendEncIds(ids []int64) {
	for _, id := range ids {
		s.EncounteredCatsId = append(s.EncounteredCatsId, int64(id))
	}
}

// BASE CURD
func (s *SelectedAnimal4Prefer) GenerateKey() int64 { // TODO 之后迁移到 model_redis 的基类去。
	s.key = variable.SnowFlake.GetId()
	s.notPassNewCat = true
	return s.key
}

func (s *SelectedAnimal4Prefer) GetKey() int64 { // TODO 同上
	return s.key
}

func (s *SelectedAnimal4Prefer) GetDataByKey(key int64) (bool, error) {
	s.key = key

	redisClient := redis_factory.GetOneRedisClient()
	defer redisClient.ReleaseOneRedisClient()

	if res, err := redisClient.String(redisClient.Execute("get", key)); err == nil {
		json.Unmarshal([]byte(res), s)
		return true, nil
	} else {
		return false, err
	}
}

func (s *SelectedAnimal4Prefer) SetDataByKey() (bool, error) {
	redisClient := redis_factory.GetOneRedisClient()
	defer redisClient.ReleaseOneRedisClient()

	data, _ := json.Marshal(s)
	if _, err := redisClient.Execute("set", s.key, data); err == nil {
		return true, nil
	} else {
		return false, err
	}
}
