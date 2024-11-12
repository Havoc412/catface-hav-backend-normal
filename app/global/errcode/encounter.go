package errcode

const (
	ErrEaLinkInstert = ErrEncounter + iota
)

func EnocunterMsgInit(m msg) {
	m[ErrEaLinkInstert] = "路遇添加成功，但关联毛茸茸失败"
}

func EncounterMsgUserInit(m msg) {
	m[ErrEaLinkInstert] = "路遇上传成功啦。但是不小心让毛茸茸跑丢了 /(ㄒoㄒ)/~~"
}
