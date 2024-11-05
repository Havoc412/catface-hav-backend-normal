package nlp

import (
	"catface/app/global/variable"
	"catface/app/utils/nlp/glm"
)

func GenerateTitle(content string) string {
	message := variable.PromptsYml.GetString("Prompt.Title") + content
	title, _ := glm.Chat(message)
	return title
}
