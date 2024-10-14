package errcode

const (
	ErrAnimalSqlFind = iota + ErrAnimal
	AnimalNoFind
)

func AnimalMsgInit(m msg) {
	m[ErrAnimalSqlFind] = "Animals 表单查询失败"
	m[AnimalNoFind] = "Animals 没有查询到符合条件的目标"
}
