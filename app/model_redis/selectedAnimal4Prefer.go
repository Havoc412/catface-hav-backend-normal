package model_redis

// INFO 辅助 animal - list - prefer 模式下的查询特化。

func CreateSelectedAnimal4Prefer() *SelectedAnimal4Prefer { // 构造函数的必要性，这样才能调用【父类】的方法。
	return &SelectedAnimal4Prefer{
		BaseModel: NewBaseModel(),
	}
}

type SelectedAnimal4Prefer struct {
	*BaseModel        `json:"-"` // TIP Go 中组合的概念。
	EncounteredCatsId []int64    `json:"encountered_cats_id"` // #1 对应第一阶段：近期路遇关联
	NewCatsId         []int64    `json:"new_cats_id"`         // #2 对应第二阶段：近期新增

	notPassNewCat bool `json:"-"`
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

/**
 * @description: 实现重写，补充子类需要增加的操作。
 * @return {*}
 */
func (s *SelectedAnimal4Prefer) GenerateKey() int64 {
	s.notPassNewCat = true
	return s.BaseModel.GenerateKey()
}

// 调用 BaseModel 的 SetDataByKey 方法
func (s *SelectedAnimal4Prefer) SetDataByKey() (bool, error) {
	return s.BaseModel.SetDataByKey(s)
}

// 调用 BaseModel 的 GetDataByKey 方法
func (s *SelectedAnimal4Prefer) GetDataByKey(key int64) (bool, error) {
	return s.BaseModel.GetDataByKey(key, s)
}
