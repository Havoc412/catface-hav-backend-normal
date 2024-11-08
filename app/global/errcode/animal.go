package errcode

const (
	ErrAnimalSqlFind = iota + ErrAnimal
	AnimalNoFind
)

func AnimalMsgInit(m msg) {
	m[ErrAnimalSqlFind] = "Animals 表单查询失败"
	m[AnimalNoFind] = "Animals 没有查询到符合条件的目标"
}

func AnimalMsgUserInit(m msg) {
	m[AnimalNoFind] = "没有更多符合此条件的毛茸茸啦，试着更换查询条件或者新增吧~"
}
