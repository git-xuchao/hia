package spark

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/urfave/cli.v1"

	"hia/cmd/spark/types"
	/*
	 *"hia/spark/ysdb"
	 */)

var routes = types.Routes{
	types.Route{"Index", "GET", "/", index},
	types.Route{"RegisterUser", "POST", "/users/:userType", registerUser},
	types.Route{"UploadVideo", "POST", "/videos/:videoID", uploadVideo},
	types.Route{"DeleteVideo", "DELETE", "/videos/:videoID", deleteVideo},
	types.Route{"PurchaseVideo", "POST", "/transaction/:videoID", purchaseVideo},
	types.Route{"PlayVideo", "GET", "/videos/:videoID", playVideo},
	types.Route{"Search", "GET", "/users/:indexType", search},
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func registerUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user types.User

	base := GetGlobalBase()
	db := base.db
	ethcli := base.ethclient

	fmt.Fprintf(w, "register, userType %s!\n", ps.ByName("userType"))
	fmt.Printf("register, userType %s!\n", ps.ByName("userType"))
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)

	if err := json.Unmarshal(body, &user); err == nil {
		fmt.Println("username:", user.UserName, ", Password:", user.Password, ", Id:", user.UserID, ", UserType:", user.UserType)
		account := ethcli.NewAccount(user.Password)
		len := len(account)
		key, _ := ethcli.GetKey(account[3 : len-1])
		fmt.Printf("key :%s", key)
		user.EthAccount = account
		user.EthKey = key
		user.EthKeyFileName, _ = ethcli.GetKeyFileName(account)

		fmt.Print("user info:\n")
		fmt.Println(user)

		err = db.UserAdd(&user)
		if err != nil {
			fmt.Printf("database add user error")
		}
	} else {
		fmt.Println("json.Unmarshal err")
		fmt.Println(err)
	}
}

func uploadVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var video types.Video

	base := GetGlobalBase()
	db := base.db
	ethcli := base.ethclient

	fmt.Fprintf(w, "uploadVideo, video name %s!\n", ps.ByName("videoID"))
	fmt.Printf("uploadVideo, video name %s!\n", ps.ByName("videoID"))
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)

	if err := json.Unmarshal(body, &video); err == nil {
		var msg ethereum.CallMsg

		fmt.Println("VideoName:", video.VideoName, ", url:", video.URL, ", UserID:", video.UserID)

		ethcli.ConstructAbi("./copyright_sol_copyright.abi")
		ethcli.SetCallMsg(&msg, "0x6e2d604754ae054e2558b38a265cb84fccb975f6", "0xa231475d813a4e642c0f98fe3167211e2e9d133d", "", "", "", nil)

		result, err := ethcli.CallContractMethod(msg, "123456", "playVideo", "alan", "http://127.0.0.1/abc.flv")
		if err != nil {
			fmt.Printf("err: %v\n", err)
		} else {
			fmt.Printf("result %v\n", result)
		}

		video.Transaction = string(result)

		db.VideoAdd(&video)

	} else {
		fmt.Println("json.Unmarshal err")
		fmt.Println(err)
	}
}

func deleteVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var video types.Video

	base := GetGlobalBase()
	db := base.db
	ethcli := base.ethclient

	fmt.Fprintf(w, "deleteVideo, video name %s!\n", ps.ByName("videoName"))
	fmt.Printf("deleteVideo, video name %s!\n", ps.ByName("videoName"))
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)

	if err := json.Unmarshal(body, &video); err == nil {
		var msg ethereum.CallMsg

		fmt.Println("VideoName:", video.VideoName, ", url:", video.URL, ", UserID:", video.UserID)

		ethcli.ConstructAbi("./copyright_sol_copyright.abi")
		ethcli.SetCallMsg(&msg, "0x6e2d604754ae054e2558b38a265cb84fccb975f6", "0xa231475d813a4e642c0f98fe3167211e2e9d133d", "", "", "", nil)

		result, err := ethcli.CallContractMethod(msg, "123456", "playVideo", "alan", "http://127.0.0.1/abc.flv")
		if err != nil {
			fmt.Printf("err: %v\n", err)
		} else {
			fmt.Printf("result %v\n", result)
		}

		db.VideoAdd(&video)
	} else {
		fmt.Println("json.Unmarshal err")
		fmt.Println(err)
	}
}

func purchaseVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var video types.Video

	base := GetGlobalBase()
	db := base.db
	ethcli := base.ethclient

	fmt.Fprintf(w, "purchaseVideo, video name %s!\n", ps.ByName("videoName"))
	fmt.Printf("purchaseVideo, video name %s!\n", ps.ByName("videoName"))
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)

	if err := json.Unmarshal(body, &video); err == nil {
		var msg ethereum.CallMsg

		fmt.Println("VideoName:", video.VideoName, ", url:", video.URL, ", UserID:", video.UserID)

		ethcli.ConstructAbi("./copyright_sol_copyright.abi")
		ethcli.SetCallMsg(&msg, "0x6e2d604754ae054e2558b38a265cb84fccb975f6", "0xa231475d813a4e642c0f98fe3167211e2e9d133d", "", "", "", nil)

		result, err := ethcli.CallContractMethod(msg, "123456", "playVideo", "alan", "http://127.0.0.1/abc.flv")
		if err != nil {
			fmt.Printf("err: %v\n", err)
		} else {
			fmt.Printf("result %v\n", result)
		}

		db.VideoAdd(&video)
	} else {
		fmt.Println("json.Unmarshal err")
		fmt.Println(err)
	}
}

func playVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var msg ethereum.CallMsg

	base := GetGlobalBase()
	ethcli := base.ethclient

	r.ParseForm()
	userID := r.Form["userID"][0]
	indexValue := r.Form["url"][0]
	videoID := ps.ByName("videoID")
	fmt.Println("userID", userID, "indexValue", indexValue, "videoID", videoID)

	ethcli.ConstructAbi("./copyright_sol_copyright.abi")
	ethcli.SetCallMsg(&msg, "0x6e2d604754ae054e2558b38a265cb84fccb975f6", "0xa231475d813a4e642c0f98fe3167211e2e9d133d", "", "", "", nil)

	result, err := ethcli.CallContractMethod(msg, "123456", "playVideo", "alan", "http://127.0.0.1/abc.flv")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("result %v\n", result)
	}
}

func search(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	/*
	 *base := GetGlobalBase()
	 *ethcli := base.ethclient
	 */

	fmt.Printf("search, indexType %s!\n", ps.ByName("indexType"))
	r.ParseForm()
	fmt.Println("userID", r.Form["userID"][0])
	fmt.Println("indexValue", r.Form["indexValue"][0])
}

func NewServer(ctx *cli.Context) error {

	router := NewHttpRouter()
	base, _ := NewBase()
	SetGlobalBase(base)
	log.Fatal(http.ListenAndServe(":8080", router))

	return nil
}
