package report

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
	"github.com/jwma/jump-jump/internal/app/utils"
)

type dailyReportWrapper struct {
	Key    string
	LinkId string
	Report *models.DailyReport
}

type Generator struct {
	isStop               chan bool
	db                   *redis.Client
	tasks                chan *models.ActiveLink
	reports              chan *dailyReportWrapper
	taskDispatchTicker   *time.Ticker
	needDispatchPastTask bool
}

func NewGenerator(rdb *redis.Client, duration time.Duration) *Generator {
	g := &Generator{
		db: rdb, taskDispatchTicker: time.NewTicker(duration),
		tasks: make(chan *models.ActiveLink, 5), reports: make(chan *dailyReportWrapper, 5),
		isStop: make(chan bool),
	}

	exists, _ := g.db.Exists(context.Background(), "dispatch_past_task").Result()
	if exists == 0 {
		g.needDispatchPastTask = true
	}

	return g
}

func (g *Generator) dispatchDailyTask() {
	now := time.Now()
	startTime := now.Add(-time.Second * 60)
	isYesterday := false

	if now.Hour() == 0 && now.Minute() <= 1 {
		isYesterday = true
		d := now.AddDate(0, 0, -1)
		startTime = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
	}

	repo := repository.GetActiveLinkRepo(g.db)
	activeLinks := repo.FindByDateRange(startTime, now)

	for _, one := range activeLinks {
		if isYesterday {
			g.tasks <- &models.ActiveLink{Id: one.Id, Time: one.Time.AddDate(0, 0, -1)}
		}
		g.tasks <- one
	}
}

func (g *Generator) dispatchPastTask() {
	if !g.needDispatchPastTask {
		return
	}

	// Get all link IDs from Redis active links or a tracking set
	// For P1, we scan through known patterns
	linkIds, _ := g.db.ZRange(context.Background(), utils.GetActiveLinkKey(), 0, -1).Result()
	st, _ := time.ParseInLocation("2006-01-02", "2020-03-01", time.Local)
	t := time.Now().AddDate(0, 0, 1)
	endTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	for st.Before(endTime) {
		for _, id := range linkIds {
			g.tasks <- &models.ActiveLink{Id: id, Time: st}
		}
		st = st.AddDate(0, 0, 1)
	}

	g.db.Set(context.Background(), "dispatch_past_task", 1, 0)
	g.needDispatchPastTask = false
}

func (g *Generator) calc(activeLink *models.ActiveLink) {
	g.reports <- CalcDailyReport(g.db, activeLink)
}

func (g *Generator) save(w *dailyReportWrapper) {
	repo := repository.GetDailyReportRepo(g.db)
	repo.Save(w.LinkId, w.Key, w.Report)
}

func (g *Generator) Start() error {
	defer g.taskDispatchTicker.Stop()

	for {
		select {
		case isStop := <-g.isStop:
			if isStop {
				return nil
			}
		case <-g.taskDispatchTicker.C:
			log.Println("ReportGenerator running...")
			go g.dispatchPastTask()
			go g.dispatchDailyTask()
		case task := <-g.tasks:
			go g.calc(task)
		case r := <-g.reports:
			go g.save(r)
		}
	}
}

func (g *Generator) Stop() error {
	g.isStop <- true
	close(g.isStop)
	close(g.tasks)
	return nil
}
