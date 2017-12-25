package spark

import (
	"fmt"
	"os"
	"sync"

	"github.com/Unknwon/goconfig"
	"gopkg.in/urfave/cli.v1"

	"hia/cmd/spark/ysdb"
	ethclient "hia/ethclient"
)

type Base struct {
	ethclient *ethclient.EthClient
	db        ysdb.YsDb
	lock      sync.RWMutex
	cfg       *goconfig.ConfigFile
	ctx       *cli.Context
}

var (
	GlobalBase *Base
)

func NewBase(ctx *cli.Context) (*Base, error) {
	var client *ethclient.EthClient
	var db ysdb.YsDb
	var ethclientAddr, keyStorePath, driverName, dataSourceName string

	configFilePath := ctx.GlobalString("config")
	fmt.Println(configFilePath)

	_, err := PathExists(configFilePath)
	if err != nil {
		fmt.Println("PathExists", err)
		return nil, err
	}

	var cfg *goconfig.ConfigFile
	cfg, err = goconfig.LoadConfigFile(configFilePath)

	ethclientAddr, err = cfg.GetValue("Ethereum", "ethNodeAddr")
	if err != nil {
		return nil, err
	}

	keyStorePath, err = cfg.GetValue("Ethereum", "keyStorePath")
	if err != nil {
		return nil, err
	}

	driverName, err = cfg.GetValue("DataBase", "driverName")
	if err != nil {
		return nil, err
	}

	dataSourceName, err = cfg.GetValue("DataBase", "dataSourceName")
	if err != nil {
		return nil, err
	}

	fmt.Println("ethNodeAddr", ethclientAddr, "keyStorePath", keyStorePath, "driverName", driverName, "dataSourceName", dataSourceName)

	client = ethclient.NewEthClient()
	client.Dial(ethclientAddr)
	client.SetKeyStoreSearchingPath(keyStorePath)
	/*
	 *db, _ = ysdb.NewABCDatabase()
	 */
	db = ysdb.NewDbMysql()
	/*
	 *db.Init(driverName, dataSourceName)
	 */
	db.Init("mysql", "root:root@tcp(192.168.31.19)/test")

	return &Base{
		ethclient: client,
		db:        db,
		cfg:       cfg,
	}, nil
}

func GetGlobalBase() *Base {
	return GlobalBase
}

func SetGlobalBase(b *Base) {
	GlobalBase = b
}

func NewSpark(ctx *cli.Context) error {
	/*
	 *fmt.Println(ctx.GlobalString("network1"))
	 *fmt.Println(ctx.String("network2"))
	 *fmt.Println(ctx.String("network3"))
	 */
	base, _ := NewBase(ctx)
	SetGlobalBase(base)

	NewServer(ctx)

	return nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return false, err
}
