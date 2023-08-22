package utils

import (
	"fmt"
	"strconv"
)

func GetRedisKeyWithPrefix(prefix string, bizId int64) string {
	return fmt.Sprintf("%v%v", prefix, strconv.FormatInt(bizId, 10))
}
