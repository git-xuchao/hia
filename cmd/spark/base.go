package spark

import (
	"sync"

	"hia/cmd/spark/ysdb"
	ethclient "hia/ethclient"
)

type Base struct {
	ethclient *ethclient.EthClient
	db        ysdb.YsDb
	lock      sync.RWMutex
}

var (
	GlobalBase *Base
)

func NewBase() (*Base, error) {
	var client *ethclient.EthClient
	var db ysdb.YsDb
	var ethclientAddr string = "http://192.168.31.52:8545"

	client = ethclient.NewEthClient()
	client.Dial(ethclientAddr)
	client.SetKeyStoreSearchingPath("/home/alan/tmp/data/node1/keystore")
	db, _ = ysdb.NewABCDatabase()

	return &Base{
		ethclient: client,
		db:        db,
	}, nil
}

func GetGlobalBase() *Base {
	return GlobalBase
}

func SetGlobalBase(b *Base) {
	GlobalBase = b
}
