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
	config.SetupConfig(pool, rdb)
}

func TestShortLinkRepository_Save(t *testing.T) {
	l := &models.ShortLink{
		Id:          "mj",
		TenantID:    "00000000-0000-0000-0000-000000000001",
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
	rs, err := repo.List("00000000-0000-0000-0000-000000000001", "mj", false, 0, 10)

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
		TenantID:    "00000000-0000-0000-0000-000000000001",
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

	rh := models.NewRequestHistory(l, "127.0.0.1", "fake user agent", "Linux")
	rhRepo := GetRequestHistoryRepo(getTestRDB(), getTestPool())
	rhRepo.Save(rh)
}

func TestRequestHistoryRepository_FindByDateRange(t *testing.T) {
	id := "testrh"
	rhRepo := GetRequestHistoryRepo(getTestRDB(), getTestPool())
	rs := rhRepo.FindByDateRange(id, time.Now().Add(-time.Hour*24), time.Now())

	if rs == nil {
		t.Error("expected non-nil result")
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

	u.TenantID = "00000000-0000-0000-0000-000000000001"
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
		TenantID:    "00000000-0000-0000-0000-000000000001",
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

	_, err := repo.FindOneByUsername("00000000-0000-0000-0000-000000000001", "")

	if err == nil {
		t.Errorf("expected error but got nil")
	}

	_, err = repo.FindOneByUsername("00000000-0000-0000-0000-000000000001", "anmuji")

	if err == nil {
		t.Errorf("expected error but got nil")
	}

	expectedUsername := "mj"
	u, err := repo.FindOneByUsername("00000000-0000-0000-0000-000000000001", expectedUsername)

	if err != nil {
		t.Error(err)
	}
	if u.Username != "mj" {
		t.Errorf("expected %s but got %s\n", expectedUsername, u.Username)
	}
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	repo := GetUserRepo(getTestPool())

	u, err := repo.FindOneByUsername("00000000-0000-0000-0000-000000000001", "mj")

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
