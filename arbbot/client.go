package arbbot

import (
	"sync"

	"github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/utils"
)

var c = new(client.Client)
var initOnce sync.Once

func GetClient() *client.Client {
	initOnce.Do(func() {
		endpoint := utils.GetEnvString("PN_URL", "http://pn:8551")
		c, _ = client.Dial(endpoint)
	})

	return c
}
