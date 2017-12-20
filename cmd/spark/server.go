package spark

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/urfave/cli.v1"

	"hia/cmd/spark/types"
	/*
	 *"hia/spark/ysdb"
	 */)

var routes = types.Routes{
	types.Route{"Index", "GET", "/", index},
	/*
	 *types.Route{"Hello", "GET", "/hello/:name", hello},
	 */
	types.Route{"RegisterUser", "POST", "/users/:usertype", registerUser},
	types.Route{"UploadVideo", "POST", "/videos/:videoName", uploadVideo},
	types.Route{"DeleteVideo", "DELETE", "/videos/:videoName", deleteVideo},
	/*
	 *types.Route{"PurchaseVideo", "POST", "/videos/:videoName", purchaseVideo},
	 */
	types.Route{"PlayVideo", "GET", "/videos/:videoName", playVideo},
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func registerUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user types.User

	base := GetGlobalBase()
	db := base.db
	ethcli := base.ethclient

	fmt.Fprintf(w, "register, usertype %s!\n", ps.ByName("usertype"))
	fmt.Printf("register, usertype %s!\n", ps.ByName("usertype"))
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)

	if err := json.Unmarshal(body, &user); err == nil {
		fmt.Println("json.Unmarshal user")
		fmt.Println(user)
		fmt.Println("username:", user.UserName, ", Password:", user.Password, ", Id:", user.ID, ", UserType:", user.UserType)
		db.UserAdd(&user)
		account := ethcli.NewAccount(user.Password)
		len := len(account)
		key, _ := ethcli.GetKey(account[3 : len-1])
		fmt.Printf("key :%s", key)
		user.EthAccount = account
		user.EthKey = key
		user.EthKeyFileName, _ = ethcli.GetKeyFileName(account)

		fmt.Print("user info:\n")
		fmt.Println(user)
	} else {
		fmt.Println("json.Unmarshal err")
		fmt.Println(err)
	}
}

func uploadVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var video types.Video

	fmt.Fprintf(w, "uploadVideo, video name %s!\n", ps.ByName("videoName"))
	fmt.Printf("uploadVideo, video name %s!\n", ps.ByName("videoName"))
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)

	if err := json.Unmarshal(body, &video); err == nil {
		fmt.Println("json.Unmarshal video")
		fmt.Println(video)
	} else {
		fmt.Println("json.Unmarshal err")
		fmt.Println(err)
	}
}

func deleteVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var video types.Video

	fmt.Fprintf(w, "deleteVideo, video name %s!\n", ps.ByName("videoName"))
	fmt.Printf("deleteVideo, video name %s!\n", ps.ByName("videoName"))
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)

	if err := json.Unmarshal(body, &video); err == nil {
		fmt.Println("json.Unmarshal video")
		fmt.Println(video)
	} else {
		fmt.Println("json.Unmarshal err")
		fmt.Println(err)
	}
}

func purchaseVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var video types.Video

	fmt.Fprintf(w, "purchaseVideo, video name %s!\n", ps.ByName("videoName"))
	fmt.Printf("purchaseVideo, video name %s!\n", ps.ByName("videoName"))
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)

	if err := json.Unmarshal(body, &video); err == nil {
		fmt.Println("json.Unmarshal video")
		fmt.Println(video)
	} else {
		fmt.Println("json.Unmarshal err")
		fmt.Println(err)
	}
}

func playVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var video types.Video

	fmt.Fprintf(w, "playVideo, video name %s!\n", ps.ByName("videoName"))
	fmt.Printf("playVideo, video name %s!\n", ps.ByName("videoName"))
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)

	if err := json.Unmarshal(body, &video); err == nil {
		fmt.Println("json.Unmarshal video")
		fmt.Println(video)
	} else {
		fmt.Println("json.Unmarshal err")
		fmt.Println(err)
	}
}

func NewServer(ctx *cli.Context) error {

	router := NewHttpRouter()
	base, _ := NewBase()
	SetGlobalBase(base)
	log.Fatal(http.ListenAndServe(":8080", router))

	return nil
}
