package spark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/julienschmidt/httprouter"

	"hia/cmd/spark/types"
)

func TestHttpRouter(t *testing.T) {
	router := httprouter.New()
	router.GET("/", index)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func TestHttpRouter2(t *testing.T) {
	router := NewHttpRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}

func TestHttpRegistUser(t *testing.T) {
	url := "http://127.0.0.1:8080/users/author"

	user := &types.User{
		UserName: "author02",
		Password: "123456",
		UserID:   12311113,
		UserType: "author",
	}

	post, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("json.marshal failed, err:", err)
		return
	}

	fmt.Printf("%s\n", string(post))

	var jsonStr = []byte(post)
	fmt.Println("jsonStr", jsonStr)
	fmt.Println("new_str", bytes.NewBuffer(jsonStr))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func TestHttpUpdateUser(t *testing.T) {
	dat, _ := ioutil.ReadFile("./ysdb/copyright_sol_copyright.abi")
	url := "http://127.0.0.1:8080/update/users"

	user := &types.User{
		UserName:        "xusd",
		Password:        "123456",
		UserID:          1111111125,
		UserType:        "author",
		Email:           "alan@sina.com",
		EthContractAddr: "0x23063382209741b9a9bf24f4fb861ffcbb8a3292",
		EthAbi:          string(dat),
	}

	/*
	 *user := &types.User{
	 *    UserName:        "author01",
	 *    Password:        "123456",
	 *    UserID:          1231111,
	 *    UserType:        "author",
	 *    Email:           "alan@sina.com",
	 *    EthContractAddr: "0x23063382209741b9a9bf24f4fb861ffcbb8a3292",
	 *    EthAbi:          string(dat),
	 *}
	 */

	post, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("json.marshal failed, err:", err)
		return
	}

	fmt.Printf("%s\n", string(post))

	var jsonStr = []byte(post)
	fmt.Println("jsonStr", jsonStr)
	fmt.Println("new_str", bytes.NewBuffer(jsonStr))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func TestHttpUploadVideo(t *testing.T) {
	url := "http://127.0.0.1:8080/videos/abcdefghiklmnop.flv"

	video := &types.Video{
		UserID: 987654325,
		/*
		 *UserName:  "alan",
		 */
		URL:       "http://127.0.0.1:8080/videos/abcdefghiklmnop.flv",
		VideoName: "abcdefgh.flv",
	}

	post, err := json.Marshal(video)
	if err != nil {
		fmt.Printf("json.marshal failed, err:", err)
		return
	}

	fmt.Printf("%s\n", string(post))

	var jsonStr = []byte(post)
	fmt.Println("jsonStr", jsonStr)
	fmt.Println("new_str", bytes.NewBuffer(jsonStr))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func TestHttpDeleteVideo(t *testing.T) {
	url := "http://127.0.0.1:8080/videos/abcdefghiklmno.flv"

	video := &types.Video{
		UserID: 987654325,
		/*
		 *UserName:  "alan",
		 */
		URL:       "http://127.0.0.1:8080/videos/abcdefghiklmno.flv",
		VideoName: "abc.flv",
	}

	post, err := json.Marshal(video)
	if err != nil {
		fmt.Printf("json.marshal failed, err:", err)
		return
	}

	fmt.Printf("%s\n", string(post))

	var jsonStr = []byte(post)
	fmt.Println("jsonStr", jsonStr)
	fmt.Println("new_str", bytes.NewBuffer(jsonStr))

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonStr))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func TestHttpPurchaceVideo(t *testing.T) {
	url := "http://127.0.0.1:8080/transaction/abcdefghiklmno.flv"

	/*
	 *video := &types.Video{
	 *    UserID: 987654326,
	 *    URL:       "http://127.0.0.1:8080/videos/abc.flv",
	 *    VideoName: "abc.flv",
	 *}
	 */
	video := &types.Video{
		UserID: 1581341302,
		/*
		 *UserName:  "alan",
		 */
		/*
		 *URL:       "http://127.0.0.1:8080/videos/abcdefg.flv",
		 */
		URL:       "http://127.0.0.1:8080/videos/abcdefghiklmno.flv",
		VideoName: "abcde.flv",
	}

	post, err := json.Marshal(video)
	if err != nil {
		fmt.Printf("json.marshal failed, err:", err)
		return
	}

	fmt.Printf("%s\n", string(post))

	var jsonStr = []byte(post)
	fmt.Println("jsonStr", jsonStr)
	fmt.Println("new_str", bytes.NewBuffer(jsonStr))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func TestHttpPlayVideo(t *testing.T) {
	url := "http://127.0.0.1:8080/videos/abcdefghiklmno.flv?userID=1581341302&url=http://127.0.0.1:8080/videos/abcdefghiklmnopk.flv"

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func TestHttpSearchUsers(t *testing.T) {
	/*
	 *url := "http://127.0.0.1:8080/record/users?userID=1581341302"
	 */
	/*
	 *url := "http://127.0.0.1:8080/record/users?count=a"
	 */
	url := "http://127.0.0.1:8080/record/users"
	/*
	 *url := "http://127.0.0.1:8080/record/users?userID=158134130234"
	 */

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	req.Header.Set("Content-Type", "application/json")
	/*
	 *req.Header.Add("Start-Time", "2017-07-27 8:46:15")
	 *req.Header.Add("End-Time", "2017-12-28 8:46:15")
	 */
	/*
	 *req.Header.Add("Start-Time", "2017-12-22 16:08:46")
	 */
	/*
	 *req.Header.Add("End-Time", "2017-12-28 8:47:15")
	 */

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func TestHttpSearchVideos(t *testing.T) {
	/*
	 *url := "http://127.0.0.1:8080/record/videos?videoID=23412454326"
	 */
	/*
	 *url := "http://127.0.0.1:8080/record/videos?indexType=uploadRecord&userID=987654325"
	 */
	/*
	 *url := "http://127.0.0.1:8080/record/videos?indexType=uploadRecord&userID=987654325"
	 */
	/*
	 *url := "http://127.0.0.1:8080/record/videos?indexType=uploadRecord&userID=987654325&count=3"
	 */
	/*
	 *url := "http://127.0.0.1:8080/record/videos?indexType=uploadRecord&userID=987654325&count=a"
	 */
	/*
	 *url := "http://127.0.0.1:8080/record/videos?indexType=videoState"
	 */
	url := "http://127.0.0.1:8080/record/videos?indexType=videoAttrib&videoID=abcdefgh.flv"

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Start-Time", "2016-07-27 8:46:15")
	req.Header.Add("End-Time", "2016-07-28 8:46:15")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func TestHttpSearchTransactions(t *testing.T) {
	/*
	 *url := "http://127.0.0.1:8080/record/transactions?videoID=23412454326"
	 */
	url := "http://127.0.0.1:8080/record/transactions?count=1"

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Start-Time", "2015-07-27 8:46:15")
	req.Header.Add("End-Time", "2017-12-28 8:46:15")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
