package errcode

const (
	ErrKnowledgeRandomList = iota + ErrKnowledge
)

func KnowledgeMsgInit(m msg) {
	m[ErrKnowledgeRandomList] = "随机获取科普知识点失败"
}

func KnowledgeMsgUserInit(m msg) {
	m[ErrKnowledgeRandomList] = "随机获取科普知识点失败"
}
