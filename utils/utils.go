package utils

import (
	"strings"
)

// EmptyString 判断字符串是否为空
func EmptyString(str string) bool {
	str = strings.TrimSpace(str)
	return strings.EqualFold(str, "")
}
