/**
* @file server.go
* @Synopsis
* @author alan lin
* @version 1.0
* @date 2017-12-24
 */
package spark

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/urfave/cli.v1"

	"hia/cmd/spark/types"
)

var routes = types.Routes{
	types.Route{"Index", "GET", "/", index},
	types.Route{"RegisterUser", "POST", "/users/:userType", registerUser},
	types.Route{"UploadVideo", "POST", "/videos/:videoID", uploadVideo},
	types.Route{"DeleteVideo", "DELETE", "/videos/:videoID", deleteVideo},
	types.Route{"PurchaseVideo", "POST", "/transaction/:videoID", purchaseVideo},
	types.Route{"PlayVideo", "GET", "/videos/:videoID", playVideo},
	types.Route{"SearchUsers", "GET", "/record/users", searchUsers},
	types.Route{"SearchVideos", "GET", "/record/videos", searchVideos},
	types.Route{"SearchTransactions", "GET", "/record/transactions", searchTransactions},
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
		video.Transaction = common.ToHex(result)
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

	videoIDStr := ps.ByName("videoID")
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

		/*
		 *userIDStr := strconv.FormatUint(video.UserID, 10)
		 *fmt.Println("userIDStr", userIDStr)
		 */

		result, err := ethcli.CallContractMethodPack(msg, resUser.Password, "deleteVideo", video.URL)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		} else {
			fmt.Printf("result %v\n", result)
		}

		video.Transaction = common.ToHex(result)
		video.VideoID = videoIDStr

		err = db.VideoDelete(&video)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
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
	videoIDStr := ps.ByName("videoID")
	body, _ := ioutil.ReadAll(r.Body)
	/*
	 *body_str := string(body)
	 *fmt.Println(body_str)
	 */

	if err := json.Unmarshal(body, &video); err == nil {
		var msg ethereum.CallMsg

		fmt.Println("VideoName:", video.VideoName, ", url:", video.URL, ", UserID:", video.UserID, "videoID", videoIDStr)
		/* search video info*/
		queryVideo.VideoID = videoIDStr
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

		result, err := ethcli.CallContractMethodPack(msg, resUser.Password, "purchaseVideo", userIDStr, video.URL)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			/*
			 *return
			 */
		} else {
			fmt.Printf("result %v\n", result)
		}

		/* add transaction info into db*/
		var transaction types.VideoTransaction
		transaction.UserID = video.UserID
		transaction.Transaction = common.ToHex(result)
		transaction.VideoID = videoIDStr

		fmt.Println("VideoTransactionAdd")
		err = db.VideoTransactionAdd(&transaction)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

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

func searchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user types.User
	var userID uint64
	var err error
	var timeStart, timeEnd time.Time
	var userIDStr string
	var userIDSlice []string
	var usersRes *[]types.User

	base := GetGlobalBase()
	db := base.db

	r.ParseForm()

	userIDSlice = r.Form["userID"]
	if len(userIDSlice) != 0 {
		userIDStr = r.Form["userID"][0]
	}

	timeStartStr := r.Header.Get("Start-Time")
	timeEndStr := r.Header.Get("End-Time")

	fmt.Println("userID", userIDStr, "timeStart", timeStartStr, "timeEnd", timeEndStr)

	if userIDStr != "" {
		userID, err = strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			return
		}
		user.UserID = userID
	}

	if timeStartStr != "" && timeEndStr != "" {
		timeStart, err = time.Parse("2006-01-02 15:04:05", timeStartStr)
		if err != nil {
			return
		}

		timeEnd, err = time.Parse("2006-01-02 15:04:05", timeEndStr)
		if err != nil {
			return
		}

		usersRes, err = db.UserQueryBetween(&user, timeStart, timeEnd)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if timeEndStr != "" {
		timeEnd, err = time.Parse("2006-01-02 15:04:05", timeEndStr)
		if err != nil {
			return
		}
		usersRes, err = db.UserQueryBefore(&user, timeEnd)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		timeStart, err = time.Parse("2006-01-02 15:04:05", timeStartStr)
		if err != nil {
			return
		}
		usersRes, err = db.UserQueryAfter(&user, timeStart)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if usersRes != nil {
		for _, u := range *usersRes {
			fmt.Println(u)
		}
	}

	fmt.Println("\n")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	body := struct {
		Data interface{} `json:"users"`
	}{
		Data: *usersRes,
	}

	if err := json.NewEncoder(w).Encode(body); err != nil {
		panic(err)
	}
}

func searchVideos(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var video types.Video
	var err error
	var userID uint64
	var timeStart, timeEnd time.Time
	var videoIDStr, userIDStr, indexTypeStr string
	var videosRes *[]types.Video

	base := GetGlobalBase()
	db := base.db

	r.ParseForm()

	videoIDSlice := r.Form["videoID"]
	if len(videoIDSlice) != 0 {
		videoIDStr = r.Form["videoID"][0]
	}

	userIDSlice := r.Form["userID"]
	if len(userIDSlice) != 0 {
		userIDStr = r.Form["userID"][0]
	}

	indexTypeSlice := r.Form["indexType"]
	if len(indexTypeSlice) != 0 {
		indexTypeStr = r.Form["indexType"][0]
	}

	timeStartStr := r.Header.Get("Start-Time")
	timeEndStr := r.Header.Get("End-Time")

	switch indexTypeStr {
	case "uploadRecord":
		if userIDStr == "" {
			fmt.Printf("searchVideos, indexType %s,  userID is not set", indexTypeStr)
			return
		}
	case "videoRanking":
		return

	case "videoAttrib":
		if videoIDStr == "" {
			fmt.Printf("searchVideos, indexType %s,  videoID is not set", indexTypeStr)
			return
		}
	case "videoState":
		if videoIDStr == "" {
			fmt.Printf("searchVideos, indexType %s,  videoID is not set", indexTypeStr)
			return
		}
	}

	fmt.Println("videoID", videoIDStr, "timeStart", timeStartStr, "timeEnd", timeEndStr)

	if videoIDStr != "" {
		video.VideoID = videoIDStr
	} else if userIDStr != "" {
		userID, err = strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			return
		}
		video.UserID = userID
	}

	if timeStartStr != "" && timeEndStr != "" {
		timeStart, err = time.Parse("2006-01-02 15:04:05", timeStartStr)
		if err != nil {
			return
		}

		timeEnd, err = time.Parse("2006-01-02 15:04:05", timeEndStr)
		if err != nil {
			return
		}
		videosRes, err = db.VideoQueryBetween(&video, timeStart, timeEnd)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if timeEndStr != "" {
		timeEnd, err = time.Parse("2006-01-02 15:04:05", timeEndStr)
		if err != nil {
			return
		}
		videosRes, err = db.VideoQueryBefore(&video, timeEnd)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		timeStart, err = time.Parse("2006-01-02 15:04:05", timeStartStr)
		if err != nil {
			return
		}
		videosRes, err = db.VideoQueryAfter(&video, timeStart)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if videosRes != nil {
		for _, u := range *videosRes {
			fmt.Println(u)
		}
	} else {
		return
	}

	fmt.Println("\n")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	switch indexTypeStr {
	case "uploadRecord":
		body := struct {
			Data interface{} `json:"uploadRecord"`
		}{
			Data: *videosRes,
		}
		if err := json.NewEncoder(w).Encode(body); err != nil {
			panic(err)
		}

	case "videoRanking":
		body := struct {
			Data interface{} `json:"videoRanking"`
		}{
			Data: *videosRes,
		}
		if err := json.NewEncoder(w).Encode(body); err != nil {
			panic(err)
		}

	case "videoAttrib":
		body := struct {
			Data interface{} `json:"videoAttrib"`
		}{
			Data: *videosRes,
		}
		if err := json.NewEncoder(w).Encode(body); err != nil {
			panic(err)
		}
	case "videoState":
		body := struct {
			Data interface{} `json:"videoState"`
		}{
			Data: *videosRes,
		}
		if err := json.NewEncoder(w).Encode(body); err != nil {
			panic(err)
		}

	}
}

func searchTransactions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var transaction types.VideoTransaction
	var err error
	var userID uint64
	var timeStart, timeEnd time.Time
	var videoIDStr, userIDStr string
	var videoIDSlice, userIDSlice []string
	var transactionsRes *[]types.VideoTransaction

	base := GetGlobalBase()
	db := base.db

	r.ParseForm()

	/*get pars*/
	videoIDSlice = r.Form["videoID"]
	if len(videoIDSlice) != 0 {
		videoIDStr = r.Form["videoID"][0]
	}

	userIDSlice = r.Form["userID"]
	if len(userIDSlice) != 0 {
		userIDStr = r.Form["userID"][0]
	}

	timeStartStr := r.Header.Get("Start-Time")
	timeEndStr := r.Header.Get("End-Time")

	fmt.Println("videoID", videoIDStr, "timeStart", timeStartStr, "timeEnd", timeEndStr)

	if videoIDStr != "" {
		transaction.VideoID = videoIDStr
	} else if userIDStr != "" {
		userID, err = strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			return
		}
		transaction.UserID = userID
	}

	if timeStartStr != "" && timeEndStr != "" {
		timeStart, err = time.Parse("2006-01-02 15:04:05", timeStartStr)
		if err != nil {
			return
		}

		timeEnd, err = time.Parse("2006-01-02 15:04:05", timeEndStr)
		if err != nil {
			return
		}

		transactionsRes, err = db.VideoTransactionQueryBetween(&transaction, timeStart, timeEnd)
		if err != nil {
			fmt.Println(err)
			return
		}

	} else if timeEndStr != "" {
		timeEnd, err = time.Parse("2006-01-02 15:04:05", timeEndStr)
		if err != nil {
			return
		}
		transactionsRes, err = db.VideoTransactionQueryBefore(&transaction, timeEnd)
		if err != nil {
			fmt.Println(err)
			return
		}

	} else {
		timeStart, err = time.Parse("2006-01-02 15:04:05", timeStartStr)
		if err != nil {
			return
		}
		transactionsRes, err = db.VideoTransactionQueryAfter(&transaction, timeStart)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if transactionsRes != nil {
		for _, u := range *transactionsRes {
			fmt.Println(u)
		}
	} else {
		fmt.Println("transactionsRes is nil")
		return
	}

	fmt.Println("\n")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	body := struct {
		Data interface{} `json:"transactions"`
	}{
		Data: *transactionsRes,
	}

	if err := json.NewEncoder(w).Encode(body); err != nil {
		panic(err)
	}
}

func NewServer(ctx *cli.Context) error {
	router := NewHttpRouter()
	log.Fatal(http.ListenAndServe(":8080", router))

	return nil
}
