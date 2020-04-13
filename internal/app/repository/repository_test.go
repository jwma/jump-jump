package repository

import (
	"github.com/go-redis/redis"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/utils"
	"testing"
	"time"
)

func getTestRDB() *redis.Client {
	return redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379", DB: 1})
}

func init() {
	// 清空测试使用的数据库，以便后续测试正常运作
	getTestRDB().FlushDB()
}

func TestShortLinkRepository_Save(t *testing.T) {
	l := &models.ShortLink{
		Id:          "mj",
		Url:         "http://anmuji.com",
		Description: "安木鸡",
		IsEnable:    true,
		CreatedBy:   "mj",
	}

	repo := GetShortLinkRepo(getTestRDB())
	err := repo.Save(l)

	if err != nil {
		t.Error(err)
	}
}

func TestShortLinkRepository_Get(t *testing.T) {
	id := "mj"
	repo := GetShortLinkRepo(getTestRDB())
	_, err := repo.Get(id)

	if err != nil {
		t.Error(err)
	}
}

func TestShortLinkRepository_Update(t *testing.T) {
	id := "mj"
	repo := GetShortLinkRepo(getTestRDB())
	l, err := repo.Get(id)

	if err != nil {
		t.Error(err)
	}

	params := &models.UpdateShortLinkParameter{
		Url:         "http://github.com/jwma",
		Description: "安木鸡的 Github",
		IsEnable:    true,
	}

	err = repo.Update(l, params)

	if err != nil {
		t.Error(err)
	}
}

func TestShortLinkRepository_List(t *testing.T) {
	repo := GetShortLinkRepo(getTestRDB())
	rs, err := repo.List(utils.GetUserShortLinksKey("mj"), 0, 10)

	if err != nil {
		t.Error(err)
	}

	expected := 1

	if rs.Total != int64(expected) {
		t.Errorf("expected %b but got %b\n", expected, rs.Total)
	}
}

func TestShortLinkRepository_Delete(t *testing.T) {
	id := "mj"
	repo := GetShortLinkRepo(getTestRDB())
	l, err := repo.Get(id)

	if err != nil {
		t.Error(err)
	}

	repo.Delete(l)
}

func TestRequestHistoryRepository_Save(t *testing.T) {
	l := &models.ShortLink{
		Id:          "testrh",
		Url:         "http://anmuji.com",
		Description: "",
		IsEnable:    true,
		CreatedBy:   "mj",
	}
	slRepo := GetShortLinkRepo(getTestRDB())
	err := slRepo.Save(l)

	if err != nil {
		t.Error(err)
	}

	rh := models.NewRequestHistory(l, "127.0.0.1", "fake user agent")
	rhRepo := GetRequestHistoryRepo(getTestRDB())
	rhRepo.Save(rh)
}

func TestRequestHistoryRepository_FindLatest(t *testing.T) {
	id := "testrh"
	rhRepo := GetRequestHistoryRepo(getTestRDB())
	rs, err := rhRepo.FindLatest(id, 10)
	expected := 1

	if err != nil {
		t.Error(err)
	}

	if rs.Total != expected {
		t.Errorf("expected %b but gog %b\n", expected, rs.Total)
	}
}

func TestRequestHistoryRepository_FindByDate(t *testing.T) {
	id := "testrh"
	rhRepo := GetRequestHistoryRepo(getTestRDB())
	rs, err := rhRepo.FindByDate(id, time.Now())
	expected := 1

	if err != nil {
		t.Error(err)
	}

	if rs.Total != expected {
		t.Errorf("expected %b but gog %b\n", expected, rs.Total)
	}
}

func TestRequestHistoryRepository_FindByDateRange(t *testing.T) {
	id := "testrh"
	rhRepo := GetRequestHistoryRepo(getTestRDB())
	start := time.Now()
	end := start.Add(time.Hour * 24)

	// 开始日期大于结束日期
	_, err := rhRepo.FindByDate(id, end, start)
	if err == nil {
		t.Errorf("expected error but got nil")
	}

	// 正常情况
	rs, err := rhRepo.FindByDate(id, start, end)
	expected := 1

	if err != nil {
		t.Error(err)
	}

	if rs.Total != expected {
		t.Errorf("expected %b but gog %b\n", expected, rs.Total)
	}
}
