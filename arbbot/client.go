package arbbot

import (
	"context"
	"math/big"
	"sync"

	"github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/utils"
)

type ArbBot struct {
	*client.Client
	Addr     common.Address
	GasPrice *big.Int
}

var arbBot = new(ArbBot)
var initOnce sync.Once

func GetClient() *ArbBot {
	initOnce.Do(func() {
		endpoint := utils.GetEnvString("PN_URL", "http://pn:8551")

		c, err := client.Dial(endpoint)
		if err != nil {
			panic(err)
		}
		myPrvKey := utils.GetEnvString("MY_PRV_KEY", "e5ae44e2ab03f3277a8c849d20cb1ad1b57781720320a09d3ec54bed4e793546")
		myAddr, err := c.ImportRawKey(context.Background(), myPrvKey, "newwhale")
		if err != nil {
			panic(err)
		}
		unlocked, err := c.UnlockAccount(context.Background(), myAddr, "newwhale", 0)
		if err != nil || !unlocked {
			panic(err)
		}

		arbBot.Client = c
		arbBot.Addr = myAddr
		arbBot.GasPrice, err = c.SuggestGasPrice(context.Background())
	})

	return arbBot
}
