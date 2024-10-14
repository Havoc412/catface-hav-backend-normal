package model

// INFO 一些基础表单的整合

type Breed struct {
	BriefModel
}

type Sterilzation struct { // TEST How to use BriefModel, the dif between Common
	Id     int64  `json:"id"`
	NameZh string `json:"name_zh"`
	NameEn string `json:"name_en"`
}

type AnmStatus struct {
	BriefModel
}

type AnmGender struct {
	BriefModel
}
