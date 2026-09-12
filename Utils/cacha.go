package Utils

import (
	"fmt"
	"strconv"
)

func CachePublishMsgName(deliverTag uint64, userID uint) string {
	return fmt.Sprintf("user:%s:tag:%s", strconv.Itoa(int(userID)), strconv.Itoa(int(deliverTag)))
}
