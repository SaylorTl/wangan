package Utils

import (
	"github.com/panjf2000/ants/v2"
	"wangxin2.0/databases"
)

var P *ants.Pool
var SubdomainPool *ants.Pool
var ImportAssetPool *ants.Pool
var ImportSyncAssetPool *ants.Pool

func InitAntPool() *ants.Pool {
	configPoolNum, _ := databases.Conf.Section("goroutinepool").Key("TOTAL_POOL_NUM").Int()
	P, _ = ants.NewPool(configPoolNum)
	return P
}

func InitSubdomainPool() *ants.Pool {
	SubdomainPool, _ = ants.NewPool(1)
	return SubdomainPool
}

func InitImportAssetPool() *ants.Pool {
	configPoolNum, _ := databases.Conf.Section("goroutinepool").Key("wangan_SYNC_POOL_NUM").Int()
	if 0 == configPoolNum {
		configPoolNum = 2
	}
	ImportAssetPool, _ = ants.NewPool(configPoolNum)
	return ImportAssetPool
}

func InitImportSyncAssetPool() *ants.Pool {
	configPoolNum, _ := databases.Conf.Section("goroutinepool").Key("wangan_SYNC_ASSET_POOL_NUM").Int()
	if 0 == configPoolNum {
		configPoolNum = 200
	}
	ImportSyncAssetPool, _ = ants.NewPool(configPoolNum)
	return ImportSyncAssetPool
}
