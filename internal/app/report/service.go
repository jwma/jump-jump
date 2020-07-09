package report

import (
	"github.com/go-redis/redis"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
	"github.com/mssola/user_agent"
	"github.com/thoas/go-funk"
	"time"
)

func CalcDailyReport(db *redis.Client, activeLink *models.ActiveLink) *dailyReportWrapper {
	date := activeLink.Time
	startTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endTime := startTime.AddDate(0, 0, 1)
	rhRepo := repository.GetRequestHistoryRepo(db)

	rhRs := rhRepo.FindByDateRange(activeLink.Id, startTime, endTime)
	pv := len(rhRs)
	ips := make([]string, 0)
	operateSystems := make(map[string]int)

	for _, rh := range rhRs {
		if !funk.ContainsString(ips, rh.IP) {
			ips = append(ips, rh.IP)
		}

		ua := user_agent.New(rh.UA)
		osInfo := ua.OSInfo()
		k := osInfo.Name

		if _, ok := operateSystems[k]; !ok {
			operateSystems[k] = 0
		}

		operateSystems[k] += 1
	}

	return &dailyReportWrapper{
		Key: startTime.Format("2006-01-02"), LinkId: activeLink.Id,
		Report: &models.DailyReport{PV: pv, UV: len(ips), OS: operateSystems},
	}
}
