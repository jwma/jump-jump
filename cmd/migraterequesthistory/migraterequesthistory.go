package main

import (
	"fmt"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/utils"
	"sync"
)

var wg sync.WaitGroup

func migrate(id string) {
	rdb := db.GetRedisClient()
	rhKey := utils.GetRequestHistoryKey(id)

	if exists, _ := rdb.Exists(rhKey).Result(); exists == 1 {
		fmt.Printf("%s 无需迁移\n", rhKey)
		wg.Done()
		return
	}

	q := fmt.Sprintf("%s:*", rhKey)

	keys, _ := rdb.Keys(q).Result()
	if len(keys) == 0 {
		wg.Done()
		return
	}

	for _, key := range keys {
		r, _ := rdb.LRange(key, 0, 99999).Result()
		rdb.RPush(rhKey, r)
	}

	rdb.Del(keys...)

	wg.Done()
}

func main() {
	rdb := db.GetRedisClient()
	c, err := rdb.ZCard(utils.GetShortLinksKey()).Result()

	if err != nil {
		panic(err)
	}

	fmt.Printf("总共有%d个短链接的访问记录可能需要迁移\n", c)
	ids, err := rdb.ZRange(utils.GetShortLinksKey(), 0, c).Result()

	if err != nil {
		panic(err)
	}

	for _, id := range ids {
		wg.Add(1)
		go migrate(id)
	}

	wg.Wait()
	fmt.Println("迁移完毕")
}
