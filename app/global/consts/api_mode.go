package consts

const (
	// TAG animal/ 接口模式
	AnimalModePrefer string = "prefer" // 根据用户行为记录优先返回【偏好目标】

	// TAG rag/chat 接口模式；  配合 yml 文件的书写习惯。
	RagChatModeKnowledge string = "Knowledge"
	RagChatModeDiary     string = "Diary"  // 查询路遇资料等
	RagChatModeDetect    string = "Detect" // 辅助 catface 的辨认功能；
)
