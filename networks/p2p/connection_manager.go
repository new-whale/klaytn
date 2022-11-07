package p2p

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/klaytn/klaytn/networks/p2p/discover"
	"github.com/klaytn/klaytn/utils"
)

type ConnectionManager struct {
	rdb      *redis.Client
	initOnce sync.Once
}

var CM = new(ConnectionManager)

func (cm *ConnectionManager) initIfNeeded() {
	cm.initOnce.Do(func() {
		addr := utils.GetEnvString("REDIS_ENDPOINT", "redis:6379")
		password := utils.GetEnvString("REDIS_PASSWORD", "")
		db := utils.GetEnvInt("REDIS_DB", 0)

		cm.rdb = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db, // use default DB
		})
	})
}

func (cm *ConnectionManager) Register(srv *BaseServer, target discover.NodeID) (bool, error) {
	cm.initIfNeeded()

	myId := discover.PubkeyID(&srv.PrivateKey.PublicKey).String()
	key := fmt.Sprintf("connection-%s", target.String())
	ctx := context.Background()

	isSet, err := cm.rdb.SetNX(ctx, key, myId, 24*time.Hour).Result()
	if err != nil {
		return false, err
	}

	if !isSet {
		id, err := cm.rdb.Get(ctx, key).Result()
		if err != nil {
			return false, err
		}
		return id == myId, nil
	}

	return true, nil
}
