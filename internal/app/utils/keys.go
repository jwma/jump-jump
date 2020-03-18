package utils

import (
	"fmt"
	"time"
)

func GetUserKey() string {
	return "users"
}

func GetShortLinkKey(id string) string {
	return fmt.Sprintf("link:%s", id)
}

func GetShortLinksKey() string {
	return "links"
}

func GetUserShortLinksKey(username string) string {
	return fmt.Sprintf("links:%s", username)
}

func GetRequestHistoryKey(linkId string, d time.Time) string {
	return fmt.Sprintf("history:%s:%s", linkId, d.Format("20060102"))
}
