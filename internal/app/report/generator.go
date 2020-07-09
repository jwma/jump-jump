package report

import (
	"github.com/go-redis/redis"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
	"github.com/jwma/jump-jump/internal/app/utils"
	"log"
	"time"
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

func NewGenerator(db *redis.Client, duration time.Duration) *Generator {
	g := &Generator{
		db: db, taskDispatchTicker: time.NewTicker(duration),
		tasks: make(chan *models.ActiveLink, 5), reports: make(chan *dailyReportWrapper, 5),
		isStop: make(chan bool),
	}

	// 设置是否需要生成过往日报的标识
	exists, _ := g.db.Exists(utils.GetDispatchPastTaskFlagKey()).Result()
	if exists == 0 {
		g.needDispatchPastTask = true
	}

	return g
}

// 日常生成日报/前一天日报
func (g *Generator) dispatchDailyTask() {
	now := time.Now()
	startTime := now.Add(-time.Second * 60) // 获取一分钟内活跃过的链接
	isYesterday := false

	// 如果当前时间区间在 00:00-00:01
	// 则将查询活跃链接的开始时间范围扩展至前一天 00:00:00
	if now.Hour() == 0 && now.Minute() <= 1 {
		isYesterday = true
		d := now.AddDate(0, 0, -1)
		startTime = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
	}

	repo := repository.GetActiveLinkRepo(g.db)
	activeLinks := repo.FindByDateRange(startTime, now)

	// 为这些链接生成/更新报表数据
	for _, one := range activeLinks {
		if isYesterday {
			g.tasks <- &models.ActiveLink{Id: one.Id, Time: one.Time.AddDate(0, 0, -1)}
		}

		g.tasks <- one
	}
}

// 生成过往日报
func (g *Generator) dispatchPastTask() {
	if !g.needDispatchPastTask {
		return
	}

	linkIds, _ := g.db.ZRange(utils.GetShortLinksKey(), 0, -1).Result()
	// 默认开始生成过往报表的开始日期为 2020-03-01，因为在那一天才加入了历史记录
	st, _ := time.ParseInLocation("2006-01-02", "2020-03-01", time.Local)
	t := time.Now().AddDate(0, 0, 1)
	endTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	for st.Before(endTime) {
		for _, id := range linkIds {
			g.tasks <- &models.ActiveLink{Id: id, Time: st}
		}

		st = st.AddDate(0, 0, 1)
	}

	g.db.Set(utils.GetDispatchPastTaskFlagKey(), 1, 0)
	g.needDispatchPastTask = false
}

// 计算报表数据
func (g *Generator) calc(activeLink *models.ActiveLink) {
	g.reports <- CalcDailyReport(g.db, activeLink)
}

// 保存报表
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

func (g *Generator) stop() error {
	g.isStop <- true
	close(g.isStop)
	close(g.tasks)
	return nil
}

func (g *Generator) Stop() error {
	return g.stop()
}
