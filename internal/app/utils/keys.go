package utils

import "fmt"

func GetShortLinkCacheKey(id string) string {
	return fmt.Sprintf("cache:link:%s", id)
}

func GetDomainCacheKey(domain string) string {
	return fmt.Sprintf("cache:domain:%s", domain)
}

func GetTenantConfigCacheKey(tenantID string) string {
	return fmt.Sprintf("cache:config:%s", tenantID)
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
