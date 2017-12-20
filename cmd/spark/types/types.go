package types

import (
	"github.com/julienschmidt/httprouter"
)

type Route struct {
	Name    string
	Method  string
	Pattern string
	Handle  httprouter.Handle
}

type Routes []Route

type User struct {
	UserName        string `json:"userName"`
	Password        string `json:"password"`
	ID              string `json:id`
	UserType        int    `json:usertype` //enum ('common','author')
	RegisterTime    uint64 `json:"registerTime"`
	UserIdCard      string `json:"userIdCard"`
	Email           string `json:"email"`
	EthAccount      string `json:"account"`
	EthKey          string `json:"secretKey"`
	EthKeyFileName  string `json:"ethKeyFileName"`
	EthContractAddr string `json:"ethrAddr"`
	EthAbi          string `json:"abi"`
	LastUpdateTime  uint64 `json:"lastUpdateTime"`
}

type Video struct {
	ID          string `json:id`
	UserName    string `json:userName`
	VideoName   string `json:videoName`
	URL         string `json:url`
	UploadTime  uint64 `json:"uploadTime"`
	Transaction string `json:"transaction"`
	Status      bool   `json:"status"`
	Plays       uint   `json:"plays"`
	Buys        uint   `json:"buys"`
}

type VideoTransaction struct {
	BuyTime       uint64 `json:"buyTime"`
	TransactionId string `json:"transactionId"`
	Url           string `json:"url"`
	UserId        uint64 `json:"userId"`
	Transaction   string `json:"transaction"`
}
