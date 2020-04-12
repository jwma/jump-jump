package repository

import (
	"github.com/go-redis/redis"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/utils"
	"testing"
)

func getTestRDB() *redis.Client {
	return redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379", DB: 1})
}

func init() {
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
