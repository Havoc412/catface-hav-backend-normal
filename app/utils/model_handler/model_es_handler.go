package model_handler

/**
 * @description: 用于处理 ES-highlight 模块分析出来的 []String.
 * @param {[]interface{}} strs
 * @return {*}
 */
func TransStringSliceToString(strs []interface{}) string {
	var result string
	for _, str := range strs {
		if s, ok := str.(string); ok {
			result += s
		}
	}
	return result
}
