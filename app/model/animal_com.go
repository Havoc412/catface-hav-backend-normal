package model

// TAG 一些基础表单的整合
type AnmSterilzation struct { // TEST How to use BriefModel, the dif between Common
	Id     int64  `json:"id"`
	NameZh string `json:"name_zh"`
	NameEn string `json:"name_en"`
}

type AnmStatus struct {
	BriefModel
	*Explain // “在校状态” 这个处理命名方式比较抽象，需要给 AI 解释一下实际含义。
}

type AnmGender struct {
	BriefModel
}

type AnmVaccination struct {
	BriefModel
}

type AnmDeworming struct {
	BriefModel
}

// TAG 带了点函数处理
func CreateAnmBreedFactory(sqlType string) *AnmBreed {
	return &AnmBreed{BriefModel: BriefModel{DB: UseDbConn(sqlType)}}
}

type AnmBreed struct {
	BriefModel
}

func (a *AnmBreed) TableName() string {
	return "anm_breeds"
}

func (a *AnmBreed) GetNameZhByEn(name_en string) string {
	var temp AnmBreed
	if result := a.DB.Where("name_en = ?", name_en).First(&temp); result.Error != nil {
		return "" // 如果查询失败，返回空字符串
	}
	return temp.NameZh
}
