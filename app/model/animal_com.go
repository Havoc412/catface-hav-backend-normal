package model

// INFO 一些基础表单的整合

type AnmBreed struct {
	BriefModel
}

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
