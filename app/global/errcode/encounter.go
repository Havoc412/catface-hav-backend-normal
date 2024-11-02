package errcode

const (
	ErrEaLinkInstert = ErrEncounter + iota
)

func EnocunterMsgInit(m msg) {
	m[ErrEaLinkInstert] = "路遇添加成功，但关联毛茸茸失败"
}