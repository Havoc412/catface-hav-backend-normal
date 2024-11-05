package nlp

import (
	"catface/app/global/variable"
	"catface/app/utils/nlp/glm"
	"catface/app/utils/yml_config/ymlconfig_interf"
)

var PromptsYml ymlconfig_interf.YmlConfigInterf

func init() {
	PromptsYml = variable.ConfigYml.Clone("rag")
}
func GenerateTitle(content string) string {
	message := PromptsYml.GetString("Prompt.Title") + content
	title, _ := glm.Chat(message)
	return title
}
