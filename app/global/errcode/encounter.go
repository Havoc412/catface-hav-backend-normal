package errcode

const (
	ErrEaLinkInstert = ErrEncounter + iota
	ErrEncounterNoData
)

func EnocunterMsgInit(m msg) {
	m[ErrEaLinkInstert] = "路遇添加成功，但关联毛茸茸失败"
	m[ErrEncounterNoData] = "没有查询到数据"
}

func EncounterMsgUserInit(m msg) {
	m[ErrEaLinkInstert] = "路遇上传成功啦。但是不小心让毛茸茸跑丢了。/(ㄒoㄒ)/~~"
	m[ErrEncounterNoData] = "没有查询到你喜欢的猫猫路遇，为他们添加吧。"
}
