package utils

import (
	"fmt"
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

func GetRequestHistoryKey(linkId string) string {
	return fmt.Sprintf("rh:%s", linkId)
}

func GetActiveLinkKey() string {
	return "activelinks"
}
