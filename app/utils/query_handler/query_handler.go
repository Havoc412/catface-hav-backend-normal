package query_handler

import (
	"strconv"
	"strings"
)

/**
 * @brie 字符串转uint8数组
 */
func StringToUint8Array(in string) []uint8 {
	var out []uint8
	_arr := strings.Split(in, ",")
	for _, c := range _arr {
		tmp, _ := strconv.ParseUint(c, 10, 8)
		out = append(out, uint8(tmp))
	}
	return out
}

func StringToStringArray(in string) []string {
	return strings.Split(in, ",")
}
