package repository

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jwma/jump-jump/internal/app/config"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/redis/go-redis/v9"
)

func getTestRDB() *redis.Client {
	return redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379", DB: 1})
}

func getTestPool() *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), "postgres://jumpjump:jumpjump@127.0.0.1:5432/jumpjump_test?sslmode=disable")
	if err != nil {
		panic(err)
	}
	return pool
}

func init() {
	rdb := getTestRDB()
	rdb.FlushDB(context.Background())

	pool := getTestPool()
	config.SetupConfig(pool)
}

func TestShortLinkRepository_Save(t *testing.T) {
	l := &models.ShortLink{
		Id:          "mj",
		Url:         "http://anmuji.com",
		Description: "安木鸡",
		IsEnable:    true,
		CreatedBy:   "mj",
	}

	repo := GetShortLinkRepo(getTestPool(), getTestRDB())
	err := repo.Save(l)

	if err != nil {
		t.Error(err)
	}
}

func TestShortLinkRepository_Get(t *testing.T) {
	id := "mj"
	repo := GetShortLinkRepo(getTestPool(), getTestRDB())
	_, err := repo.Get(id)

	if err != nil {
		t.Error(err)
	}
}

func TestShortLinkRepository_Update(t *testing.T) {
	id := "mj"
	repo := GetShortLinkRepo(getTestPool(), getTestRDB())
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
	repo := GetShortLinkRepo(getTestPool(), getTestRDB())
	rs, err := repo.List("mj", false, 0, 10)

	if err != nil {
		t.Error(err)
	}

	expected := 1

	if rs.Total != int64(expected) {
		t.Errorf("expected %d but got %d\n", expected, rs.Total)
	}
}

func TestShortLinkRepository_Delete(t *testing.T) {
	id := "mj"
	repo := GetShortLinkRepo(getTestPool(), getTestRDB())
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
	slRepo := GetShortLinkRepo(getTestPool(), getTestRDB())
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
		t.Errorf("expected %d but got %d\n", expected, rs.Total)
	}
}

func TestRequestHistoryRepository_FindByDateRange(t *testing.T) {
	id := "testrh"
	rhRepo := GetRequestHistoryRepo(getTestRDB())
	rs := rhRepo.FindByDateRange(id, time.Now().Add(-time.Second*10), time.Now())
	expected := 1

	if len(rs) != expected {
		t.Errorf("expected %d but got %d\n", expected, len(rs))
	}
}

func TestUserRepository_Save(t *testing.T) {
	repo := GetUserRepo(getTestPool())

	u := &models.User{
		Username:    "",
		Role:        0,
		RawPassword: "",
	}

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
	repo := GetUserRepo(getTestPool())

	_, err := repo.FindOneByUsername("")

	if err == nil {
		t.Errorf("expected error but got nil")
	}

	_, err = repo.FindOneByUsername("anmuji")

	if err == nil {
		t.Errorf("expected error but got nil")
	}

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
	repo := GetUserRepo(getTestPool())

	u, err := repo.FindOneByUsername("mj")

	if err != nil {
		t.Error(err)
	}

	err = repo.UpdatePassword(u)

	if err == nil {
		t.Errorf("expected error but got nil")
	}

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
