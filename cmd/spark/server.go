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
	types.Route{"updateUser", "POST", "/update/users", updateUser},
	types.Route{"UploadVideo", "POST", "/videos/:videoID", uploadVideo},
	types.Route{"DeleteVideo", "DELETE", "/videos/:videoID", deleteVideo},
	types.Route{"PurchaseVideo", "POST", "/transaction/:videoID", purchaseVideo},
	types.Route{"PlayVideo", "GET", "/videos/:videoID", playVideo},
	types.Route{"SearchUsers", "GET", "/record/users", searchUsers},
	types.Route{"SearchVideos", "GET", "/record/videos", searchVideos},
	types.Route{"SearchTransactions", "GET", "/record/transactions", searchTransactions},
}

func checkRegisterUserParamater(user types.User) error {
	fmt.Println("username:", user.UserName, ", Password:", user.Password, ", Id:", user.UserID, ", UserType:", user.UserType)

	if user.UserName == "" || user.Password == "" || user.UserID == 0 || user.UserType == "" {
		return errors.New("checking  user register paramaters error")
	}

	if user.UserType != "common" && user.UserType != "author" {
		return errors.New("checking  user register paramaters, user type error")
	}

	return nil
}

func checkParamaterWhenUploadingVideo(video types.Video) error {
	fmt.Println("VideoName:", video.VideoName, ", url:", video.URL, ", UserID:", video.UserID)

	if video.VideoName == "" || video.URL == "" || video.UserID == 0 {
		return errors.New("checking  video uploading paramaters error")
	}

	return nil
}

func checkParamaterWhenDeletingingVideo(video types.Video) error {
	fmt.Println("VideoID:", video.VideoID, ", url:", video.URL, ", UserID:", video.UserID)
	if video.VideoID == "" || video.UserID == 0 {
		return errors.New("checking  video deleteing paramaters error")
	}
	return nil
}

func checkParamaterWhenPurchasingingVideo(video types.Video) error {
	fmt.Println("VideoName:", video.VideoName, ", url:", video.URL, ", UserID:", video.UserID, "videoID", video.VideoID)
	if video.VideoID == "" || video.UserID == 0 {
		return errors.New("checking purchasing video paramaters error")
	}
	return nil
}

func checkParamaterWhenplayingVideo(video types.Video) error {
	fmt.Println("VideoName:", video.VideoName, ", url:", video.URL, ", UserID:", video.UserID, "videoID", video.VideoID)
	if video.URL == "" || video.URL == "" || video.UserID == 0 {
		return errors.New("checking playing video paramaters error")
	}
	return nil
}

func checkParamaterWhenSearchingUsers(userID, timeStart, timeEnd, count string) error {
	fmt.Println("userID", userID, "timeStart", timeStart, "timeEnd", timeEnd, "count", count)
	/*
	 *if userID == "" && timeStart == "" && timeEnd == "" {
	 *    return errors.New("checking searching users paramaters error")
	 *}
	 */
	return nil
}

func checkParamaterWhenSearchingTransactions(videoID, userID, timeStart, timeEnd string) error {
	fmt.Println("videoID", videoID, "userID", userID, "timeStart", timeStart, "timeEnd", timeEnd)
	/*
	 *if "videoID" == "" && userID == "" && timeStart == "" && timeEnd == "" {
	 *    return errors.New("checking  searching transactions paramaters error")
	 *}
	 */
	return nil
}

func checkParamaterWhenSearchingVideos(indexType, videoID, userID, timeStart, timeEnd string) error {

	fmt.Println("indexType", indexType, "videoID", videoID, "userID", userID, "timeStart", timeStart, "timeEnd", timeEnd)

	if indexType == "" || (videoID == "" && userID == "" && timeStart == "" && timeEnd == "") {
		return errors.New("checking  searching videos paramaters error")
	}

	switch indexType {
	case "uploadRecord":
		if userID == "" {
			return errors.New("search Videos uploadRecord, userID is not set")
		}
	case "videoRanking":
		/*
		 *return errors.New("search Videos videoRanking, not support yet")
		 */

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

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func registerUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user types.User

	base := GetGlobalBase()
	db := base.db
	ethcli := base.ethclient

	fmt.Printf("register, userType %s!\n", ps.ByName("userType"))
	userTypeStr := ps.ByName("userType")
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)

	if err := json.Unmarshal(body, &user); err == nil {
		user.UserType = userTypeStr

		/*1.check user register infos*/
		err = checkRegisterUserParamater(user)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		/*2.search if user exists*/
		var usersRes *[]types.User
		var userSearch types.User
		userSearch.UserID = user.UserID
		userSearch.UserName = user.UserName
		usersRes, err = db.UserQuery(&userSearch, "")
		if usersRes != nil {
			fmt.Println("user exists")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		/*3.new a account in ethereum*/
		account := ethcli.NewAccount(user.Password)
		len := len(account)
		var key string
		key, err = ethcli.GetKey(account[3 : len-1])

		/*4.write infos to db*/
		user.EthAccount = account[1 : len-1]
		user.EthKey = key
		user.EthKeyFileName, _ = ethcli.GetKeyFileName(account)
		fmt.Print("user info:\n")
		fmt.Println(user)
		err = db.UserAdd(&user)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			/*
			 *if err := json.NewEncoder(w).Encode(&user); err != nil {
			 *    fmt.Println(err)
			 *    w.WriteHeader(http.StatusInternalServerError)
			 *    return
			 *}
			 */
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "{}")
		}

	} else {
		fmt.Println("json.Unmarshal err")
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func updateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user types.User

	base := GetGlobalBase()
	db := base.db

	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)

	if err := json.Unmarshal(body, &user); err == nil {
		/*update user*/

		err = db.UserUpdate(&user)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
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

	/*
	 *fmt.Fprintf(w, "uploadVideo, video name %s!\n", ps.ByName("videoID"))
	 */
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
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		/*2.search if video has uploaded*/
		var videoRes *[]types.Video
		var videoSearch types.Video
		videoSearch.VideoID = video.VideoID
		videoRes, err = db.VideoQuery(&videoSearch, "")
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if videoRes != nil {
			if len(*videoRes) != 0 {
				fmt.Println("video has uploaded")
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		/*3.search author user info*/
		user.UserID = video.UserID
		resUser, err = db.UserQuerySimple(&user)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if resUser.UserType != "author" {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println("author user info", resUser)

		/*4.call contract method*/
		ethcli.ConstructAbi2(resUser.EthAbi)
		fmt.Println("EthAccount", resUser.EthAccount, "EthContractAddr", resUser.EthContractAddr)
		ethcli.SetCallMsg(&msg, resUser.EthAccount, resUser.EthContractAddr, "", "", "", nil)
		result, err := ethcli.CallContractMethodPack(msg, resUser.Password, "uploadVideo", video.URL)
		if err != nil {
			if err.Error() == "video has been uploaded,please not repeat upload;" {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			} else {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			fmt.Printf("result %v\n", result)
		}

		/*5.write db*/
		video.Transaction = common.ToHex(result)
		fmt.Println("adfds")
		fmt.Println("video", video)
		fmt.Println("VideoAdd")
		err = db.VideoAdd(&video)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			fmt.Println("sdfasdfasdf")
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "{}")
		}

	} else {
		fmt.Println("json.Unmarshal err")
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func deleteVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var video types.Video
	var user, resUser types.User

	base := GetGlobalBase()
	db := base.db
	ethcli := base.ethclient

	fmt.Printf("deleteVideo, video name %s!\n", ps.ByName("videoName"))

	videoIDStr := ps.ByName("videoID")
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)

	if err := json.Unmarshal(body, &video); err == nil {
		var msg ethereum.CallMsg

		/*1.check check paramater when deleting video*/
		video.VideoID = videoIDStr
		err = checkParamaterWhenDeletingingVideo(video)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		/*2.search if video exists*/
		var videoRes *[]types.Video
		var videoSearch types.Video
		videoSearch.VideoID = videoIDStr
		videoRes, err = db.VideoQuery(&videoSearch, "")
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if videoRes == nil {
			fmt.Println("video not exists")
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			video.URL = (*videoRes)[0].URL
		}

		/*3.search author user info*/
		user.UserID = video.UserID
		resUser, err = db.UserQuerySimple(&user)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			fmt.Println(resUser)
		}
		if resUser.UserType != "author" {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		/*4.call contract*/
		ethcli.ConstructAbi2(resUser.EthAbi)
		ethcli.SetCallMsg(&msg, resUser.EthAccount, resUser.EthContractAddr, "", "", "", nil)
		result, err := ethcli.CallContractMethodPack(msg, resUser.Password, "deleteVideo", video.URL)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			fmt.Printf("result %v\n", result)
		}

		/*5.write db*/
		video.VideoID = videoIDStr
		err = db.VideoDelete(&video)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		} else {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "{}")
		}
	} else {
		fmt.Println("json.Unmarshal err")
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func purchaseVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var video, queryVideo, resVideo types.Video
	var user, resUser types.User

	base := GetGlobalBase()
	db := base.db
	ethcli := base.ethclient

	fmt.Printf("purchaseVideo, video name %s!\n", ps.ByName("videoID"))
	videoIDStr := ps.ByName("videoID")
	body, _ := ioutil.ReadAll(r.Body)
	/*
	 *body_str := string(body)
	 *fmt.Println(body_str)
	 */

	if err := json.Unmarshal(body, &video); err == nil {
		var msg ethereum.CallMsg

		/*1.check paramaters when purchasing video*/
		video.VideoID = videoIDStr
		err = checkParamaterWhenPurchasingingVideo(video)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		/*2.search if user exists*/
		var usersRes *[]types.User
		var userSearch types.User
		userSearch.UserID = video.UserID
		usersRes, err = db.UserQuery(&userSearch, "")
		if usersRes == nil {
			fmt.Println("user not exists")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		/*2.search if this video has purchased*/
		var transactionSearch types.VideoTransaction
		var transactionsRes *[]types.VideoTransaction
		transactionSearch.UserID = video.UserID
		transactionSearch.VideoID = videoIDStr
		transactionsRes, err = db.VideoTransactionQuery(&transactionSearch, "")
		if transactionsRes != nil {
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			} else if len(*transactionsRes) != 0 {
				fmt.Println("user has purchased this video")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		/*3.search video info*/
		queryVideo.VideoID = videoIDStr
		resVideo, err = db.VideoQuerySimple(&queryVideo)
		if err != nil {
			if err.Error() == "result is null" {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			} else {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			fmt.Println("search video info", resVideo)
			video.URL = resVideo.URL
		}

		/*4.search author user info*/
		user.UserID = resVideo.UserID
		resUser, err = db.UserQuerySimple(&user)
		if err != nil {
			if err.Error() == "result is null" {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			} else {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
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

		/*5.call contract method*/
		ethcli.ConstructAbi2(resUser.EthAbi)
		ethcli.SetCallMsg(&msg, resUser.EthAccount, resUser.EthContractAddr, "", "", "", nil)

		result, err := ethcli.CallContractMethodPack(msg, resUser.Password, "purchaseVideo", userIDStr, video.URL)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			fmt.Printf("result %v\n", result)
		}

		/*6.add transaction info into db*/
		var transaction types.VideoTransaction
		transaction.UserID = video.UserID
		transaction.Transaction = common.ToHex(result)
		transaction.VideoID = videoIDStr
		fmt.Println("VideoTransactionAdd")
		err = db.VideoTransactionAdd(&transaction)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "{}")
		}

	} else {
		fmt.Println("json.Unmarshal err")
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func playVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var msg ethereum.CallMsg
	var queryVideo, resVideo, checkVideo types.Video
	var user, resUser types.User
	var err error
	var videoIDStr, userIDStr string

	base := GetGlobalBase()
	db := base.db
	ethcli := base.ethclient

	/*parse params */
	r.ParseForm()
	userIDStr = r.Form["userID"][0]
	url := r.Form["url"][0]
	videoIDStr = ps.ByName("videoID")
	fmt.Println("userID", userIDStr, "url", url, "videoID", videoIDStr)

	/*1.check paramater when playing video*/
	checkVideo.URL = url
	checkVideo.UserID, _ = strconv.ParseUint(userIDStr, 10, 64)
	err = checkParamaterWhenplayingVideo(checkVideo)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	/*2.search if this video has purchased*/
	var transactionSearch types.VideoTransaction
	var transactionsRes *[]types.VideoTransaction
	transactionSearch.UserID, err = strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	transactionSearch.VideoID = videoIDStr
	transactionsRes, err = db.VideoTransactionQuery(&transactionSearch, "")
	if transactionsRes == nil {
		fmt.Println("user has not purchased this video")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	/*3.search video info*/
	queryVideo.URL = url
	resVideo, err = db.VideoQuerySimple(&queryVideo)
	if err != nil {
		if err.Error() == "result is null" {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		fmt.Println("search play video info:", resVideo)
		if resVideo.URL != url || resVideo.VideoID != videoIDStr {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	/*4.search author user info*/
	user.UserID = resVideo.UserID
	fmt.Println("userIDStr", userIDStr)
	resUser, err = db.UserQuerySimple(&user)
	if err != nil {
		if err.Error() == "result is null" {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		fmt.Println(resUser)
	}

	/*5.call contract method*/
	ethcli.ConstructAbi2(resUser.EthAbi)
	ethcli.SetCallMsg(&msg, resUser.EthAccount, resUser.EthContractAddr, "", "", "", nil)
	result, err := ethcli.CallContractMethodOnly(msg, nil, "playVideo", userIDStr, url)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		fmt.Printf("result %v\n", result)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"userName\":%s, \"videoName\":%s, \"url\":%s}", resUser.UserName, resVideo.VideoName, url)
	}
}

func searchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user types.User
	var userID uint64
	var count int = 10
	var err error
	var timeStart, timeEnd time.Time
	var userIDStr, countStr string
	var userIDSlice, countSlice []string
	var usersRes *[]types.User

	base := GetGlobalBase()
	db := base.db

	r.ParseForm()

	/*1.get pars*/
	userIDSlice = r.Form["userID"]
	if len(userIDSlice) != 0 {
		userIDStr = r.Form["userID"][0]
	}

	countSlice = r.Form["count"]
	if len(countSlice) != 0 {
		countStr = r.Form["count"][0]
	}

	timeStartStr := r.Header.Get("Start-Time")
	timeEndStr := r.Header.Get("End-Time")

	/*2.check check paramater when search users*/
	err = checkParamaterWhenSearchingUsers(userIDStr, timeStartStr, timeEndStr, countStr)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	/*3.convert pars*/
	if userIDStr != "" {
		userID, err = strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		user.UserID = userID
	}

	if countStr != "" {
		/*
		 *count, err = strconv.ParseUint(countStr, 10, 64)
		 */
		count, err = strconv.Atoi(countStr)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println("count", count)
	}

	if timeStartStr != "" {
		timeStart, err = time.Parse("2006-01-02 15:04:05", timeStartStr)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if timeEndStr != "" {
		timeEnd, err = time.Parse("2006-01-02 15:04:05", timeEndStr)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	/*4.search db*/
	if timeStartStr != "" && timeEndStr != "" {
		usersRes, err = db.UserQueryBetween(&user, timeStart, timeEnd, 0, count)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if timeEndStr != "" {
		usersRes, err = db.UserQueryBefore(&user, timeEnd, 0, count)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if timeStartStr != "" {
		usersRes, err = db.UserQueryAfter(&user, timeStart, 0, count)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		usersRes, err = db.UserQuery(&user, "")
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if usersRes == nil {
		fmt.Println("\n")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		return
	} else if usersRes != nil {
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

	if err = json.NewEncoder(w).Encode(body); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func searchVideos(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var video types.Video
	var err error
	var count int = 10
	var userID uint64
	var timeStart, timeEnd time.Time
	var videoIDStr, userIDStr, indexTypeStr, countStr string
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

	countSlice := r.Form["count"]
	if len(countSlice) != 0 {
		countStr = r.Form["count"][0]
	}

	timeStartStr := r.Header.Get("Start-Time")
	timeEndStr := r.Header.Get("End-Time")

	/*2.check check paramater when searching video */
	err = checkParamaterWhenSearchingVideos(indexTypeStr, videoIDStr, userIDStr, timeStartStr, timeEndStr)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	/*3.convert datas*/
	if videoIDStr != "" {
		video.VideoID = videoIDStr
	} else if userIDStr != "" {
		userID, err = strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		video.UserID = userID
	}

	if countStr != "" {
		count, err = strconv.Atoi(countStr)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println("count", count)
	}

	if timeStartStr != "" {
		timeStart, err = time.Parse("2006-01-02 15:04:05", timeStartStr)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if timeEndStr != "" {
		timeEnd, err = time.Parse("2006-01-02 15:04:05", timeEndStr)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	/*4.search db*/
	if indexTypeStr == "videoRanking" {
		sqlSlip := fmt.Sprintf("1 ORDER BY plays limit %d", count)
		videosRes, err = db.VideoQuery(&video, sqlSlip)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		if timeStartStr != "" && timeEndStr != "" {
			videosRes, err = db.VideoQueryBetween(&video, timeStart, timeEnd, 0, count)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else if timeEndStr != "" {
			videosRes, err = db.VideoQueryBefore(&video, timeEnd, 0, count)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else if timeStartStr != "" {
			videosRes, err = db.VideoQueryAfter(&video, timeStart, 0, count)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
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
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case "videoRanking":
		body := struct {
			Data interface{} `json:"videoRanking"`
		}{
			Data: *videosRes,
		}
		if err := json.NewEncoder(w).Encode(body); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case "videoAttrib":
		body := struct {
			Data interface{} `json:"videoAttrib"`
		}{
			Data: *videosRes,
		}
		if err := json.NewEncoder(w).Encode(body); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "videoState":
		body := struct {
			Data interface{} `json:"videoState"`
		}{
			Data: *videosRes,
		}
		if err := json.NewEncoder(w).Encode(body); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
}

func searchTransactions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var transaction types.VideoTransaction
	var err error
	var count int = 10
	var userID uint64
	var timeStart, timeEnd time.Time
	var videoIDStr, userIDStr, countStr string
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

	countSlice := r.Form["count"]
	if len(countSlice) != 0 {
		countStr = r.Form["count"][0]
	}

	timeStartStr := r.Header.Get("Start-Time")
	timeEndStr := r.Header.Get("End-Time")

	/*2.check check paramaters when searching transactions */
	err = checkParamaterWhenSearchingTransactions(videoIDStr, userIDStr, timeStartStr, timeEndStr)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	/*3.convert datas*/
	if videoIDStr != "" {
		transaction.VideoID = videoIDStr
	} else if userIDStr != "" {
		userID, err = strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		transaction.UserID = userID
	}

	if countStr != "" {
		count, err = strconv.Atoi(countStr)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println("count", count)
	}

	if timeStartStr != "" {
		timeStart, err = time.Parse("2006-01-02 15:04:05", timeStartStr)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if timeEndStr != "" {
		timeEnd, err = time.Parse("2006-01-02 15:04:05", timeEndStr)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	/*4.search db*/
	if timeStartStr != "" && timeEndStr != "" {
		transactionsRes, err = db.VideoTransactionQueryBetween(&transaction, timeStart, timeEnd, 0, count)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if timeEndStr != "" {
		transactionsRes, err = db.VideoTransactionQueryBefore(&transaction, timeEnd, 0, count)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if timeStartStr != "" {
		transactionsRes, err = db.VideoTransactionQueryAfter(&transaction, timeStart, 0, count)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		transactionsRes, err = db.VideoTransactionQuery(&transaction, "")
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
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
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func NewServer(ctx *cli.Context) error {
	router := NewHttpRouter()
	log.Fatal(http.ListenAndServe(":8080", router))

	return nil
}
