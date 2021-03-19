package repository

import (
	"github.com/go-redis/redis"
	"github.com/jwma/jump-jump/internal/app/config"
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

	config.SetupConfig(getTestRDB())
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

	params := &models.UpdateShortLinkAPIRequest{
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
		t.Errorf("expected %b but got %b\n", expected, rs.Total)
	}
}

func TestRequestHistoryRepository_FindByDateRange(t *testing.T) {
	id := "testrh"
	rhRepo := GetRequestHistoryRepo(getTestRDB())
	rs := rhRepo.FindByDateRange(id, time.Now().Add(-time.Second*10), time.Now())
	expected := 1

	if len(rs) != expected {
		t.Errorf("expected %b but got %b\n", expected, len(rs))
	}
}

func TestUserRepository_Save(t *testing.T) {
	repo := GetUserRepo(getTestRDB())

	u := &models.User{
		Username:    "",
		Role:        0,
		RawPassword: "",
	}

	// 测试保存不符合要求的用户数据
	err := repo.Save(u)

	if err == nil {
		t.Errorf("expected error but got nil")
	}

	u.Username = "mj"
	u.RawPassword = "123456"
	err = repo.Save(u)

	if err == nil {
		t.Errorf("expected error but got nil")
	}

	u.Role = models.RoleUser
	err = repo.Save(u)

	if err != nil {
		t.Error(err)
	}

	// 尝试使用已存在的用户名创建用户
	u2 := &models.User{
		Username:    "mj",
		Role:        models.RoleUser,
		RawPassword: "abcdefg",
	}
	err = repo.Save(u2)

	if err == nil {
		t.Errorf("expected error but got nil")
	}
}

func TestUserRepository_FindOneByUsername(t *testing.T) {
	repo := GetUserRepo(getTestRDB())

	// 测试 username 空字符
	_, err := repo.FindOneByUsername("")

	if err == nil {
		t.Errorf("expected error but got nil")
	}

	// 测试查找不存在的用户名
	_, err = repo.FindOneByUsername("anmuji")

	if err == nil {
		t.Errorf("expected error but got nil")
	}

	// 正常查找
	expectedUsername := "mj"
	u, err := repo.FindOneByUsername(expectedUsername)

	if err != nil {
		t.Error(err)
	}
	if u.Username != "mj" {
		t.Errorf("expected %s but got %s\n", expectedUsername, u.Username)
	}
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	repo := GetUserRepo(getTestRDB())

	u, err := repo.FindOneByUsername("mj")

	if err != nil {
		t.Error(err)
	}

	//测试更新密码为空字符串
	err = repo.UpdatePassword(u)

	if err == nil {
		t.Errorf("expected error but got nil")
	}

	//测试正常更新密码
	u.RawPassword = "opqrst"
	err = repo.UpdatePassword(u)

	if err != nil {
		t.Error(err)
	}
}

func TestActiveLinkRepository_Save(t *testing.T) {
	repo := GetActiveLinkRepo(getTestRDB())
	repo.Save("a")
	repo.Save("b")
	repo.Save("c")
}

func TestActiveLinkRepository_FindByDateRange(t *testing.T) {
	repo := GetActiveLinkRepo(getTestRDB())
	activeLinks := repo.FindByDateRange(time.Now().Add(-time.Minute), time.Now())
	expected := 3

	if len(activeLinks) != expected {
		t.Errorf("expected %d but got %d", expected, len(activeLinks))
	}
}

func TestDailyReportRepository_Save(t *testing.T) {
	repo := GetDailyReportRepo(getTestRDB())
	repo.Save("fake", "2020-01-01", &models.DailyReport{
		PV: 1,
		UV: 1,
		OS: map[string]int{"Mac OS X": 1},
	})
}

func TestDailyReportRepository_FindRecent(t *testing.T) {
	repo := GetDailyReportRepo(getTestRDB())

	sampleKey := time.Now().Format("2006-01-02")
	sample := &models.DailyReport{
		PV: 1,
		UV: 1,
		OS: map[string]int{"Mac OS X": 1},
	}
	repo.Save("fake", sampleKey, sample)

	reports := repo.FindRecent("fake", 3)
	expected := 3

	if len(reports) != expected {
		t.Errorf("expected %d but got %d", expected, len(reports))
	}

	if reports[expected-1].Date != sampleKey {
		t.Errorf("expected %s but got %s", sampleKey, reports[expected-1].Date)
	}

	if reports[expected-1].Report.PV != sample.PV {
		t.Errorf("expected %d but got %d", sample.PV, reports[expected-1].Report.PV)
	}

	if reports[expected-1].Report.UV != sample.UV {
		t.Errorf("expected %d but got %d", sample.PV, reports[expected-1].Report.UV)
	}
}
