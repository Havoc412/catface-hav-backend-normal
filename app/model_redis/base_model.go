package model_redis

import (
	"catface/app/global/variable"
	"catface/app/utils/redis_factory"
	"encoding/json"
)

func NewBaseModel() *BaseModel {
	return &BaseModel{}
}

type BaseModel struct {
	key int64 `json:"-"`
}

func (b *BaseModel) GenerateKey() int64 {
	b.key = variable.SnowFlake.GetId()
	return b.key
}

func (b *BaseModel) GetKey() int64 {
	return b.key
}

// SetDataByKey 将子类的数据保存到 Redis 中
func (b *BaseModel) SetDataByKey(data interface{}) (bool, error) {
	redisClient := redis_factory.GetOneRedisClient()
	defer redisClient.ReleaseOneRedisClient()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	if _, err := redisClient.Execute("set", b.key, jsonData); err == nil {
		return true, nil
	} else {
		return false, err
	}
}

// GetDataByKey 从 Redis 中获取数据并解码到子类中
func (b *BaseModel) GetDataByKey(key int64, data interface{}) (bool, error) {
	b.key = key // 顺便保存 valid 中拿到的 key 值。

	redisClient := redis_factory.GetOneRedisClient()
	defer redisClient.ReleaseOneRedisClient()

	if res, err := redisClient.String(redisClient.Execute("get", key)); err == nil {
		if err := json.Unmarshal([]byte(res), data); err == nil {
			return true, nil
		} else {
			return false, err
		}
	} else {
		return false, err
	}
}
