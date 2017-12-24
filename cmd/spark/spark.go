package spark

import (
	"fmt"
	"sync"

	"gopkg.in/urfave/cli.v1"
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
	/*
	 *db, _ = ysdb.NewABCDatabase()
	 */
	db = ysdb.NewDbMysql()

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

func NewSpark(ctx *cli.Context) error {
	fmt.Println(ctx.GlobalString("network1"))
	fmt.Println(ctx.String("network2"))
	fmt.Println(ctx.String("network3"))
	NewServer(ctx)
	base, _ := NewBase()
	SetGlobalBase(base)

	return nil
}
