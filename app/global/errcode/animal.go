package errcode

const (
	ErrAnimalSqlFind = iota + ErrAnimal
)

func AnimalMsgInit(m msg) {
	m[ErrAnimalSqlFind] = "Animals 表单查询失败"
}
