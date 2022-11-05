package p2p

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/klaytn/klaytn/networks/p2p/discover"
)

type ConnectionManager struct {
	rdb      *redis.Client
	initOnce sync.Once
}

var CM ConnectionManager

func (cm ConnectionManager) InitIfNeeded() {
	cm.initOnce.Do(func() {
		// addr := os.Getenv("REDIS_ENDPOINT")

		rdb := redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "",
			DB:       0, // use default DB
		})

		cm.rdb = rdb
	})
}

func (cm ConnectionManager) Register(srv *BaseServer, target discover.NodeID) (bool, error) {
	myId := discover.PubkeyID(&srv.PrivateKey.PublicKey).String()
	key := fmt.Sprintf("connection-%s", target.String())
	ctx := context.Background()

	if _, err := cm.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		cnt, err := p.Exists(ctx, key).Result()
		if err != nil {
			return err
		}

		if cnt == 0 {
			p.Set(ctx, key, myId, 24*time.Hour)
		}

		return nil
	}); err != nil {
		return false, err
	}

	id, err := cm.rdb.Get(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return id == myId, nil
}
