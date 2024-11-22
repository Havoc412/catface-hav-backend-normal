package errcode

const (
	ErrAnimalSqlFind = iota + ErrAnimal
	AnimalNoFind

	// TAG
	CatFaceFail
	// CatFaceNoFind // INFO 交给前端判断更合适
)

func AnimalMsgInit(m msg) {
	m[ErrAnimalSqlFind] = "Animals 表单查询失败"
	m[AnimalNoFind] = "Animals 没有查询到符合条件的目标"
	m[CatFaceFail] = "猫脸识别失败"
}

func AnimalMsgUserInit(m msg) {
	m[AnimalNoFind] = "没有更多符合此条件的毛茸茸啦，试着更换查询条件或者新增吧~"
	m[CatFaceFail] = "猫脸识别失败，请重新尝试~ 😿"
}
