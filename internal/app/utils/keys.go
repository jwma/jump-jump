package utils

import (
	"fmt"
	"time"
)

func GetRequestHistoryKey(linkId string, d time.Time) string {
	return fmt.Sprintf("history:%s:%s", linkId, d.Format("20060102"))
}
