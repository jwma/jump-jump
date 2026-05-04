package utils

import "fmt"

func GetShortLinkCacheKey(id string) string {
	return fmt.Sprintf("cache:link:%s", id)
}

func GetRequestHistoryKey(linkId string) string {
	return fmt.Sprintf("rh:%s", linkId)
}

func GetActiveLinkKey() string {
	return "activelinks"
}

func GetDailyReportKey(linkId string) string {
	return fmt.Sprintf("dr:%s", linkId)
}
