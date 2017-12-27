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
	"errors"
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

func checkRegisterUserParamater(user types.User) error {
	fmt.Println("username:", user.UserName, ", Password:", user.Password, ", Id:", user.UserID, ", UserType:", user.UserType)

	if user.UserName == "" && user.Password == "" && user.UserID == 0 && user.UserType == "" {
		return errors.New("checking  user register paramaters error")
	}
	return nil
}

func registerUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user types.User

	base := GetGlobalBase()
	db := base.db
	ethcli := base.ethclient

	fmt.Fprintf(w, "register, userType %s!\n", ps.ByName("userType"))
	fmt.Printf("register, userType %s!\n", ps.ByName("userType"))
	userTypeStr := ps.ByName("userType")
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)

	if err := json.Unmarshal(body, &user); err == nil {
		if user.UserType == "" {
			user.UserType = userTypeStr
		}

		/*1.check user register infos*/
		err = checkRegisterUserParamater(user)
		if err != nil {
			fmt.Println("registerUser StatusBadRequest")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		/*2.new a account in ethereum*/
		account := ethcli.NewAccount(user.Password)
		len := len(account)
		key, _ := ethcli.GetKey(account[3 : len-1])
		fmt.Printf("key :%s", key)

		/*3.write infos to db*/
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

func checkParamaterWhenUploadingVideo(video types.Video) error {
	fmt.Println("VideoName:", video.VideoName, ", url:", video.URL, ", UserID:", video.UserID)

	if video.VideoName == "" && video.URL == "" && video.UserID == 0 {
		return errors.New("checking  video uploading paramaters error")
	}

	return nil
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

		/*1.check check paramater when uploading video*/
		video.VideoID = videoID
		err := checkParamaterWhenUploadingVideo(video)
		if err != nil {
			fmt.Println("uploadVideo StatusBadRequest")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		/*2.search user info*/
		user.UserID = video.UserID
		resUser, err = db.UserQuerySimple(&user)
		fmt.Println(resUser)

		/*3.call contract method*/
		ethcli.ConstructAbi2(resUser.EthAbi)
		ethcli.SetCallMsg(&msg, resUser.EthAccount, resUser.EthContractAddr, "", "", "", nil)
		result, err := ethcli.CallContractMethodPack(msg, resUser.Password, "uploadVideo", video.URL)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		} else {
			fmt.Printf("result %v\n", result)
		}

		/*4.write db*/
		video.Transaction = common.ToHex(result)
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

func checkParamaterWhenDeletingingVideo(video types.Video) error {
	fmt.Println("VideoName:", video.VideoName, ", url:", video.URL, ", UserID:", video.UserID)
	if video.VideoName == "" && video.URL == "" && video.UserID == 0 {
		return errors.New("checking  video uploading paramaters error")
	}
	return nil
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

		/*1.check check paramater when deleting video*/
		err = checkParamaterWhenDeletingingVideo(video)
		if err != nil {
			fmt.Println("deleteVideo StatusBadRequest")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		/*2.search user info*/
		user.UserID = video.UserID
		resUser, err = db.UserQuerySimple(&user)
		if err != nil {
			return
		} else {
			fmt.Println(resUser)
		}

		/*3.call contract*/
		ethcli.ConstructAbi2(resUser.EthAbi)
		ethcli.SetCallMsg(&msg, resUser.EthAccount, resUser.EthContractAddr, "", "", "", nil)
		result, err := ethcli.CallContractMethodPack(msg, resUser.Password, "deleteVideo", video.URL)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		} else {
			fmt.Printf("result %v\n", result)
		}

		/*4.write db*/
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

func checkParamaterWhenPurchasingingVideo(video types.Video) error {
	fmt.Println("VideoName:", video.VideoName, ", url:", video.URL, ", UserID:", video.UserID, "videoID", video.VideoID)
	if video.VideoID == "" && video.URL == "" && video.UserID == 0 {
		return errors.New("checking purchasing video paramaters error")
	}
	return nil
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

		/*1.check check paramater when purchasing video*/
		err = checkParamaterWhenPurchasingingVideo(video)
		if err != nil {
			fmt.Println("purchaseVideo StatusBadRequest")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		/*2.search video info*/
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

		/*3.search author user info*/
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

		/*4.call contract method*/
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

		/*5.add transaction info into db*/
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

func checkParamaterWhenplayingVideo(video types.Video) error {
	fmt.Println("VideoName:", video.VideoName, ", url:", video.URL, ", UserID:", video.UserID, "videoID", video.VideoID)
	if video.VideoID == "" && video.URL == "" && video.UserID == 0 {
		return errors.New("checking playing video paramaters error")
	}
	return nil
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

	/*1.check check paramater when playing video*/
	queryVideo.URL = url
	err = checkParamaterWhenplayingVideo(queryVideo)
	if err != nil {
		fmt.Println("playVideo StatusBadRequest")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	/*2.search video info*/
	queryVideo.URL = url
	resVideo, err = db.VideoQuerySimple(&queryVideo)
	if err != nil {
		return
	} else {
		fmt.Println("search play video info:", resVideo)
	}

	/*3.search author user info*/
	user.UserID = resVideo.UserID
	fmt.Println("userIDStr", userID)
	resUser, err = db.UserQuerySimple(&user)
	if err != nil {
		return
	} else {
		fmt.Println(resUser)
	}

	/*4.call contract method*/
	ethcli.ConstructAbi2(resUser.EthAbi)
	ethcli.SetCallMsg(&msg, resUser.EthAccount, resUser.EthContractAddr, "", "", "", nil)
	result, err := ethcli.CallContractMethodOnly(msg, nil, "playVideo", userID, url)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("result %v\n", result)
	}
}

func checkParamaterWhenSearchingUsers(userID, timeStart, timeEnd string) error {
	fmt.Println("userID", userID, "timeStart", timeStart, "timeEnd", timeEnd)
	if userID == "" || timeStart == "" || timeEnd == "" {
		return errors.New("checking  searching users paramaters error")
	}
	return nil
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

	/*1.get pars*/
	userIDSlice = r.Form["userID"]
	if len(userIDSlice) != 0 {
		userIDStr = r.Form["userID"][0]
	}

	timeStartStr := r.Header.Get("Start-Time")
	timeEndStr := r.Header.Get("End-Time")

	/*2.check check paramater when search users*/
	err = checkParamaterWhenSearchingUsers("", "", "")
	if err != nil {
		fmt.Println("searchUsers StatusBadRequest")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	/*3.convert pars*/
	if userIDStr != "" {
		userID, err = strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			return
		}
		user.UserID = userID
	}
	if timeStartStr != "" {
		timeStart, err = time.Parse("2006-01-02 15:04:05", timeStartStr)
		if err != nil {
			return
		}
	}
	if timeEndStr != "" {
		timeEnd, err = time.Parse("2006-01-02 15:04:05", timeEndStr)
		if err != nil {
			return
		}
	}

	/*4.search db*/
	if timeStartStr != "" && timeEndStr != "" {
		usersRes, err = db.UserQueryBetween(&user, timeStart, timeEnd)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if timeEndStr != "" {
		usersRes, err = db.UserQueryBefore(&user, timeEnd)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
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

	/*5.write back*/
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

func checkParamaterWhenSearchingVideos(indexType, videoID, userID, timeStart, timeEnd string) error {

	fmt.Println("indexType", indexType, "videoID", videoID, "userID", userID, "timeStart", timeStart, "timeEnd", timeEnd)

	if indexType == "" && (videoID == "" || userID == "" || timeStart == "" || timeEnd == "") {
		return errors.New("checking  searching videos paramaters error")
	}

	switch indexType {
	case "uploadRecord":
		if userID == "" {
			return errors.New("search Videos uploadRecord, userID is not set")
		}
	case "videoRanking":
		return errors.New("search Videos videoRanking, not support yet")

	case "videoAttrib":
		if videoID == "" {
			return errors.New("search Videos videoAttrib, videoID is not set")
		}
	case "videoState":
		if videoID == "" {
			return errors.New("search Videos videoState, videoID is not set")
		}
	}
	return nil
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

	/*1.get pars*/
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

	/*2.check check paramater when searching video */
	err = checkParamaterWhenSearchingVideos(indexTypeStr, videoIDStr, userIDStr, timeStartStr, timeEndStr)
	if err != nil {
		fmt.Println("playVideo StatusBadRequest")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	/*3.convert datas*/
	if videoIDStr != "" {
		video.VideoID = videoIDStr
	} else if userIDStr != "" {
		userID, err = strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			return
		}
		video.UserID = userID
	}

	if timeStartStr != "" {
		timeStart, err = time.Parse("2006-01-02 15:04:05", timeStartStr)
		if err != nil {
			return
		}
	}
	if timeEndStr != "" {
		timeEnd, err = time.Parse("2006-01-02 15:04:05", timeEndStr)
		if err != nil {
			return
		}
	}

	/*4.search db*/
	if timeStartStr != "" && timeEndStr != "" {
		videosRes, err = db.VideoQueryBetween(&video, timeStart, timeEnd)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if timeEndStr != "" {
		videosRes, err = db.VideoQueryBefore(&video, timeEnd)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
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

	/*5.write back*/
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

func checkParamaterWhenSearchingTransactions(videoID, userID, timeStart, timeEnd string) error {
	fmt.Println("videoID", videoID, "userID", userID, "timeStart", timeStart, "timeEnd", timeEnd)
	if "videoID" == "" || userID == "" || timeStart == "" || timeEnd == "" {
		return errors.New("checking  searching transactions paramaters error")
	}
	return nil
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

	/*1.get pars*/
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

	/*2.check check paramaters when searching transactions */
	err = checkParamaterWhenSearchingTransactions(videoIDStr, userIDStr, timeStartStr, timeEndStr)
	if err != nil {
		fmt.Println("searchTransactions StatusBadRequest")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	/*3.convert datas*/
	if videoIDStr != "" {
		transaction.VideoID = videoIDStr
	} else if userIDStr != "" {
		userID, err = strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			return
		}
		transaction.UserID = userID
	}

	if timeStartStr != "" {
		timeStart, err = time.Parse("2006-01-02 15:04:05", timeStartStr)
		if err != nil {
			return
		}
	}

	if timeEndStr != "" {
		timeEnd, err = time.Parse("2006-01-02 15:04:05", timeEndStr)
		if err != nil {
			return
		}
	}

	/*4.search db*/
	if timeStartStr != "" && timeEndStr != "" {
		transactionsRes, err = db.VideoTransactionQueryBetween(&transaction, timeStart, timeEnd)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if timeEndStr != "" {
		transactionsRes, err = db.VideoTransactionQueryBefore(&transaction, timeEnd)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
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

	/*5.write back*/
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
