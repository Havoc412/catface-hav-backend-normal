package model_redis

// INFO 辅助 animal - list - prefer 模式下的查询特化。

type SelectedAnimal4Prefer struct {
	Key               int64   `json:"-"`                   // redis 的 key 值
	EncounteredCatsId []int64 `json:"encountered_cats_id"` // #1 对应第一阶段：近期路遇关联
	NewCatsId         []int64 `json:"new_cats_id"`         // #2 对应第二阶段：近期新增
}

func (s *SelectedAnimal4Prefer) PassNew() bool {
	return s.Key != 0
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
