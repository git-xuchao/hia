package myhttprouter

import (
	"bytes"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestHttpRouter(t *testing.T) {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func TestHttpRouter2(t *testing.T) {
	router := NewHttpRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}

func TestHttpClient(t *testing.T) {
	url := "http://127.0.0.1:8080/users/custom"
	usrId := "alan"
	pwd := "pwd1234"
	post := "{\"UserId\":\"" + usrId + "\",\"Password\":\"" + pwd + "\"}"

	fmt.Println(url, "post", post)

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
