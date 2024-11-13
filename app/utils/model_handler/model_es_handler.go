package model_handler

func TransStringSliceToString(strs []interface{}) string {
	var result string
	for _, str := range strs {
		if s, ok := str.(string); ok {
			result += s
		}
	}
	return result
}
