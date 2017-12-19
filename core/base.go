package core

import (
	"sync"

	"hia/core/ysdb"
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
	var ethclientAddr string = "http://127.0.0.1:8001"

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
