package myhttprouter

import (
	/*
	 *"encoding/json"
	 */
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

type VedioInfo struct {
	ID        string `json:id`
	UserName  string `json:userName`
	VedioName string `json:vedioName`
	URL       string `json:url`
}
