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

func GetDailyReportKey(linkId string) string {
	return fmt.Sprintf("dr:%s", linkId)
}

func GetDispatchPastTaskFlagKey() string {
	return "dispatch_past_task"
}

func GetConfigKey() string {
	return "j2config"
}

func GetLandingHostsConfigKey() string {
	return "landingHosts"
}

func GetIdLengthConfigKey() string {
	return "idLength"
}

func GetIdMinimumLengthConfigKey() string {
	return "idMinimumLength"
}

func GetIdMaximumLengthConfigKey() string {
	return "idMaximumLength"
}

func GetShortLinkNotFoundConfigKey() string {
	return "shortLinkNotFoundConfig"
}
