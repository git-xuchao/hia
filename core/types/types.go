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
	UserName string `json:"userName"`
	Password string `json:"password"`
	ID       string `json:id`
	UserType int    `json:id`
}

type Video struct {
	ID        string `json:id`
	UserName  string `json:userName`
	VideoName string `json:videoName`
	URL       string `json:url`
}
