package core

import (
	"math/big"
)

type Vedio struct {
	URL       string
	Name      string
	PlayCount int
	BuyCount  int
	Status    bool
}

type User struct {
	UserID          big.Int
	Password        string
	Type            big.Int
	Name            string
	EthAcount       string
	EthPrivateKey   string
	EthContractAddr string
	EthContractAbi  string
}
