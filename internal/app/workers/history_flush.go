package workers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/utils"
	"github.com/redis/go-redis/v9"
)

const (
	flushBatchSize = 100
	flushInterval  = 5 * time.Second
)

type HistoryFlushWorker struct {
	rdb *redis.Client
	db  *pgxpool.Pool
}

func NewHistoryFlushWorker(rdb *redis.Client, db *pgxpool.Pool) *HistoryFlushWorker {
	return &HistoryFlushWorker{rdb: rdb, db: db}
}

func (w *HistoryFlushWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Final flush before shutdown
			w.flush()
			return
		case <-ticker.C:
			w.flush()
		}
	}
}

func (w *HistoryFlushWorker) flush() {
	for {
		members, err := w.rdb.ZPopMin(context.Background(), utils.RequestHistoryBufferKey, flushBatchSize).Result()
		if err != nil {
			if err != redis.Nil {
				log.Printf("history flush: redis error: %v", err)
			}
			return
		}
		if len(members) == 0 {
			return
		}

		if err := w.batchInsert(members); err != nil {
			log.Printf("history flush: batch insert failed: %v", err)
			// Re-push failed items back to buffer
			for _, z := range members {
				data, _ := json.Marshal(z.Member)
				w.rdb.ZAdd(context.Background(), utils.RequestHistoryBufferKey, redis.Z{
					Score:  z.Score,
					Member: string(data),
				})
			}
			return
		}
	}
}

func (w *HistoryFlushWorker) batchInsert(members []redis.Z) error {
	batch := &pgx.Batch{}
	for _, z := range members {
		raw, ok := z.Member.(string)
		if !ok {
			continue
		}
		rh := &models.RequestHistory{}
		if err := json.Unmarshal([]byte(raw), rh); err != nil {
			log.Printf("history flush: unmarshal error: %v", err)
			continue
		}
		t := rh.Time
		if t.IsZero() {
			t = time.Unix(int64(z.Score), 0)
		}
		batch.Queue(
			`INSERT INTO request_histories (short_link_id, tenant_id, url, ip, ua, created_at)
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			rh.ShortLinkID, rh.TenantID, rh.Url, rh.IP, rh.UA, t,
		)
	}

	br := w.db.SendBatch(context.Background(), batch)
	defer br.Close()

	for range members {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}
	return nil
}
