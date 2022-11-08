package utils

import (
	"os"
	"strconv"
	"strings"

	"github.com/klaytn/klaytn/common"
)

func GetEnvString(key string, defaultValue string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	} else {
		return defaultValue
	}
}

func GetEnvStrings(key string, defaultValue []string) []string {
	if val, exists := os.LookupEnv(key); exists {
		return strings.Split(val, ",")
	} else {
		return defaultValue
	}
}

func GetEnvAddresses(key string, defaultValue []string) []common.Address {
	addrStrs := GetEnvStrings(key, defaultValue)
	addrs := make([]common.Address, len(addrStrs))
	for i, addrStr := range addrStrs {
		addrs[i] = common.HexToAddress(addrStr)
	}
	return addrs
}

func GetEnvInt(key string, defaultValue int) int {
	if val, exists := os.LookupEnv(key); exists {
		v, err := strconv.Atoi(val)
		if err != nil {
			panic(err)
		}
		return v
	} else {
		return defaultValue
	}
}
