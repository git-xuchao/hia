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
	url := "http://127.0.0.1:8080/users/custom"

	user := &types.User{
		UserName: "tim",
		Password: "123456",
		UserID:   1581341302,
		UserType: "common",
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

func TestHttpUploadVideo(t *testing.T) {
	url := "http://127.0.0.1:8080/videos/abc.flv"

	video := &types.Video{
		UserID: 987654325,
		/*
		 *UserName:  "alan",
		 */
		URL:       "http://127.0.0.1:8080/videos/abc.flv",
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
	url := "http://127.0.0.1:8080/videos/abc.flv"

	video := &types.Video{
		UserID: 1581341302,
		/*
		 *UserName:  "alan",
		 */
		URL:       "http://127.0.0.1:8080/videos/abc.flv",
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
	url := "http://127.0.0.1:8080/transaction/abc.flv"

	video := &types.Video{
		UserID: 987654326,
		/*
		 *UserName:  "alan",
		 */
		URL:       "http://127.0.0.1:8080/videos/abc.flv",
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
	url := "http://127.0.0.1:8080/videos/adsfasdfasdfa?userID=987654326&url=http://127.0.0.1:8080/videos/abc.flv"

	video := &types.Video{
		UserID: 987654326,
		/*
		 *UserName:  "alan",
		 */
		URL:       "http://127.0.0.1:8080/videos/abc.flv",
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

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
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

func TestHttpSearch(t *testing.T) {
	/*
	 *url := "http://127.0.0.1:8080/record/users?userID=987654326&timeStart=2016-07-27 8:46:15&timeEnd=2016-07-28 8:46:15"
	 */
	url := "http://127.0.0.1:8080/record/users?userID=987654326&timeStart=2016-07-27&timeEnd=2016-07-28"

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
	 *url := "http://127.0.0.1:8080/record/users?userID=987654326&timeStart=2016-07-27 8:46:15&timeEnd=2016-07-28 8:46:15"
	 */
	url := "http://127.0.0.1:8080/record/users?userID=987654326"

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	// req.Header.Set("X-Custom-Header", "myvalue")
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

func TestHttpSearchVideos(t *testing.T) {
	url := "http://127.0.0.1:8080/record/videos?videoID=23412454326"

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	// req.Header.Set("X-Custom-Header", "myvalue")
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
	url := "http://127.0.0.1:8080/record/transactions?videoID=23412454326"

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Start-Time", "2016-07-27 8:46:15")
	/*
	 *req.Header.Add("End-Time", "2016-07-28 8:46:15")
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
