package spark

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

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
	types.Route{"Search", "GET", "/record/users", searchUser},
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
	var user, resUser types.User

	base := GetGlobalBase()
	db := base.db
	ethcli := base.ethclient

	fmt.Fprintf(w, "uploadVideo, video name %s!\n", ps.ByName("videoID"))
	fmt.Printf("uploadVideo, video name %s!\n", ps.ByName("videoID"))

	videoID := ps.ByName("videoID")
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)

	if err := json.Unmarshal(body, &video); err == nil {
		var msg ethereum.CallMsg

		fmt.Println("VideoName:", video.VideoName, ", url:", video.URL, ", UserID:", video.UserID)

		/* search user info*/
		user.UserID = video.UserID

		resUser, err = db.UserQuerySimple(&user)

		fmt.Println(resUser)

		/*
		 *ethcli.ConstructAbi("./copyright_sol_copyright.abi")
		 */
		/*
		 *dat, _ := ioutil.ReadFile("./copyright_sol_copyright.abi")
		 *ethcli.ConstructAbi2(string(dat))
		 */

		ethcli.ConstructAbi2(resUser.EthAbi)

		ethcli.SetCallMsg(&msg, resUser.EthAccount, resUser.EthContractAddr, "", "", "", nil)

		result, err := ethcli.CallContractMethodPack(msg, resUser.Password, "uploadVideo", video.URL)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		} else {
			fmt.Printf("result %v\n", result)
		}

		//..........
		/*
		 *video.Transaction = string(result)
		 */
		video.Transaction = "asdafasdfdasfdasfsdf"
		fmt.Println("video.Transaction", video.Transaction)
		video.VideoID = videoID
		video.URL = video.URL

		fmt.Println("adfds")
		fmt.Println("video", video)
		fmt.Println("VideoAdd")
		err = db.VideoAdd(&video)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		} else {
			fmt.Println("sdfasdfasdf")
		}

	} else {
		fmt.Println("json.Unmarshal err")
		fmt.Println(err)
	}
}

func deleteVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var video types.Video
	var user, resUser types.User

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

		/* search user info*/
		user.UserID = video.UserID
		resUser, err = db.UserQuerySimple(&user)
		if err != nil {
			return
		} else {
			fmt.Println(resUser)
		}

		ethcli.ConstructAbi2(resUser.EthAbi)
		ethcli.SetCallMsg(&msg, resUser.EthAccount, resUser.EthContractAddr, "", "", "", nil)

		result, err := ethcli.CallContractMethodPack(msg, resUser.Password, "deleteVideo", video.UserID, video.URL)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		} else {
			fmt.Printf("result %v\n", result)
		}
		video.Transaction = string(result)

		db.VideoUpdate(&video)
	} else {
		fmt.Println("json.Unmarshal err")
		fmt.Println(err)
	}
}

func purchaseVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var video, queryVideo, resVideo types.Video
	var user, resUser types.User

	base := GetGlobalBase()
	db := base.db
	ethcli := base.ethclient

	fmt.Fprintf(w, "purchaseVideo, video name %s!\n", ps.ByName("videoName"))
	fmt.Printf("purchaseVideo, video name %s!\n", ps.ByName("videoID"))
	videoID := ps.ByName("videoID")
	body, _ := ioutil.ReadAll(r.Body)
	/*
	 *body_str := string(body)
	 *fmt.Println(body_str)
	 */

	if err := json.Unmarshal(body, &video); err == nil {
		var msg ethereum.CallMsg

		fmt.Println("VideoName:", video.VideoName, ", url:", video.URL, ", UserID:", video.UserID, "videoID", videoID)
		/* search video info*/
		queryVideo.VideoID = videoID
		/*
		 *queryVideo.URL = video.URL
		 */
		resVideo, err = db.VideoQuerySimple(&queryVideo)
		if err != nil {
			return
		} else {
			fmt.Println("search video info", resVideo)
		}

		/* search author user info*/
		user.UserID = resVideo.UserID
		resUser, err = db.UserQuerySimple(&user)
		if err != nil {
			return
		} else {
			fmt.Println("search author user info", resUser)
		}

		fmt.Println("asdfasdf")
		fmt.Println("EthAbi", resUser.EthAbi)
		fmt.Println("EthAccount", resUser.EthAccount)
		fmt.Println("EthContractAddr", resUser.EthContractAddr)
		fmt.Println("Password", resUser.Password)
		fmt.Println("UserID", video.UserID)
		fmt.Println("URL", video.URL)
		userIDStr := strconv.FormatUint(video.UserID, 10)
		fmt.Println("userIDStr", userIDStr)

		ethcli.ConstructAbi2(resUser.EthAbi)
		ethcli.SetCallMsg(&msg, resUser.EthAccount, resUser.EthContractAddr, "", "", "", nil)

		/*
		 *result, err := ethcli.CallContractMethodPack(msg, resUser.Password, "purchaseVideo", video.UserID, video.URL)
		 */
		result, err := ethcli.CallContractMethodPack(msg, resUser.Password, "purchaseVideo", userIDStr, video.URL)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		} else {
			fmt.Printf("result %v\n", result)
		}

		/* add transaction info into db*/
		var transaction types.VideoTransaction
		transaction.UserID = video.UserID
		/*
		 *transaction.URL = video.URL
		 */
		transaction.Transaction = string(result)

		fmt.Println("VideoTransactionAdd")
		db.VideoTransactionAdd(&transaction)

	} else {
		fmt.Println("json.Unmarshal err")
		fmt.Println(err)
	}
}

func playVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var msg ethereum.CallMsg
	var queryVideo, resVideo types.Video
	var user, resUser types.User
	var err error

	base := GetGlobalBase()
	db := base.db
	ethcli := base.ethclient

	/*parse params */
	r.ParseForm()
	userID := r.Form["userID"][0]
	url := r.Form["url"][0]
	videoID := ps.ByName("videoID")
	fmt.Println("userID", userID, "url", url, "videoID", videoID)

	/*search video info*/
	queryVideo.URL = url
	resVideo, err = db.VideoQuerySimple(&queryVideo)
	if err != nil {
		return
	} else {
		fmt.Println("search play video info:", resVideo)
	}

	/*search author user info*/
	user.UserID = resVideo.UserID
	fmt.Println("userIDStr", userID)
	resUser, err = db.UserQuerySimple(&user)
	if err != nil {
		return
	} else {
		fmt.Println(resUser)
	}

	/*call contract method*/
	ethcli.ConstructAbi2(resUser.EthAbi)
	ethcli.SetCallMsg(&msg, resUser.EthAccount, resUser.EthContractAddr, "", "", "", nil)
	result, err := ethcli.CallContractMethodOnly(msg, nil, "playVideo", userID, url)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("result %v\n", result)
	}
}

func searchUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user, resUser types.User

	base := GetGlobalBase()
	db := base.db
	/*
	 *ethcli := base.ethclient
	 */

	/*
	 *fmt.Printf("search, indexType %s!\n", ps.ByName("indexType"))
	 */
	r.ParseForm()
	fmt.Println("userID", r.Form["userID"][0])
	userIDStr := r.Form["userID"][0]
	userID, _ := strconv.ParseUint(userIDStr, 10, 64)
	fmt.Println(userID)

	/* search user info*/
	user.UserID = userID

	resUser, _ = db.UserQuerySimple(&user)

	fmt.Println(resUser)
	/*
	 *fmt.Println("indexValue", r.Form["indexValue"][0])
	 */
}

func NewServer(ctx *cli.Context) error {

	router := NewHttpRouter()
	base, _ := NewBase()
	SetGlobalBase(base)
	log.Fatal(http.ListenAndServe(":8080", router))

	return nil
}
