package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/utils"
	"sync"
)

var wg sync.WaitGroup

func getRequestHistoryKeyV120(id string) string {
	return fmt.Sprintf("history:%s", id)
}

func handleV120Migrate(id string) {
	rdb := db.GetRedisClient()
	rhKey := getRequestHistoryKeyV120(id)

	if exists, _ := rdb.Exists(rhKey).Result(); exists == 1 {
		fmt.Printf("%s 无需迁移\n", rhKey)
		wg.Done()
		return
	}

	keys, _ := rdb.Keys(fmt.Sprintf("%s:*", rhKey)).Result()

	if len(keys) == 0 {
		wg.Done()
		return
	}

	for _, key := range keys {
		r, _ := rdb.LRange(key, 0, -1).Result()
		rdb.RPush(rhKey, r)
	}

	rdb.Del(keys...)
	wg.Done()
}

func getShortLinkIds() []string {
	rdb := db.GetRedisClient()
	c, err := rdb.ZCard(utils.GetShortLinksKey()).Result()

	if err != nil {
		panic(err)
	}

	ids, err := rdb.ZRange(utils.GetShortLinksKey(), 0, c).Result()

	if err != nil {
		panic(err)
	}

	return ids
}

func startV120Migration() {
	ids := getShortLinkIds()
	fmt.Printf("[V1.2.0] 总共有 %d 个短链接的访问记录可能需要迁移\n", len(ids))

	for _, id := range ids {
		wg.Add(1)
		go handleV120Migrate(id)
	}

	wg.Wait()
	fmt.Println("[V1.2.0] 迁移完毕")
}

func handleV130Migrate(id string) {
	rdb := db.GetRedisClient()
	key := utils.GetRequestHistoryKey(id)
	v120key := getRequestHistoryKeyV120(id)
	r, _ := rdb.LRange(v120key, 0, -1).Result()

	for _, rhStr := range r {
		rh := &models.RequestHistory{}
		_ = json.Unmarshal([]byte(rhStr), rh)
		rh.Id = utils.RandStringRunes(6)
		rdb.ZAdd(key, redis.Z{
			Score:  float64(rh.Time.Unix()),
			Member: rh,
		})
	}

	rdb.Del(v120key)
	wg.Done()
}

func startV130Migration() {
	ids := getShortLinkIds()
	fmt.Printf("[V1.3.0] 总共有 %d 个短链接的访问记录可能需要迁移\n", len(ids))

	for _, id := range ids {
		wg.Add(1)
		go handleV130Migrate(id)
	}

	wg.Wait()
	fmt.Println("[V1.3.0] 迁移完毕")
}

func main() {
	startV120Migration()
	fmt.Println()
	startV130Migration()
}
