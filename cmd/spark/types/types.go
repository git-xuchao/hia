package types

import (
	"time"

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
	UserID          uint64    `json:"userID"`
	UserName        string    `json:"userName"`
	Password        string    `json:"password"`
	UserType        string    `json:"userType"` //enum ('common','author')
	UserIdCard      string    `json:"userIdCard"`
	Email           string    `json:"email"`
	EthAccount      string    `json:"ethAccount"`
	EthKey          string    `json:"ethKey"`
	EthKeyFileName  string    `json:"ethKeyFileName"`
	EthContractAddr string    `json:"ethContractAddr"`
	EthAbi          string    `json:"ethAbi"`
	RegisterTime    time.Time `json:"registerTime"`
	LastUpdateTime  time.Time `json:"lastUpdateTime"`
}

type Video struct {
	UserID      uint64    `json:"userID"`
	VideoID     string    `json:"videoID"`
	VideoName   string    `json:"videoName"`
	URL         string    `json:"url"`
	UploadTime  time.Time `json:"uploadTime"`
	Transaction string    `json:"transaction"`
	Status      *bool     `json:"status"`
	Plays       *uint     `json:"plays"`
	Buys        *uint     `json:"buys"`
}

type VideoTransaction struct {
	BuyTime       time.Time `json:"buyTime"`
	TransactionId string    `json:"transactionId"`
	VideoID       string    `json:"videoID"`
	UserID        uint64    `json:"userId"`
	Transaction   string    `json:"transaction"`
}
