package myhttprouter

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
)

var routes = Routes{
	Route{"Index", "GET", "/", index},
	Route{"Hello", "GET", "/hello/:name", hello},
	Route{"RegisterUser", "POST", "/users/:usertype", registerUser},
	Route{"UploadVedio", "PUT", "/vedios/:vedioName", uploadVedio},
	Route{"DeleteVedio", "DELETE", "/vedios/:vedioName", deleteVedio},
	Route{"PurchaseVedio", "POST", "/vedios/:vedioName", purchaseVedio},
	Route{"PlayVedio", "GET", "/vedios/:vedioName", playVedio},
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func registerUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user User

	fmt.Fprintf(w, "register, usertype %s!\n", ps.ByName("usertype"))
	fmt.Printf("register, usertype %s!\n", ps.ByName("usertype"))
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)

	if err := json.Unmarshal(body, &user); err == nil {
		fmt.Println("json.Unmarshal user")
		fmt.Println(user)
		fmt.Println("username:", user.UserName, ", Password:", user.Password, ", Id:", user.ID, ", UserType:", user.UserType)
	} else {
		fmt.Println("json.Unmarshal err")
		fmt.Println(err)
	}
}

func uploadVedio(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var vedioInfo VedioInfo

	fmt.Fprintf(w, "uploadVedio, vedio name %s!\n", ps.ByName("vedioName"))
	fmt.Printf("uploadVedio, vedio name %s!\n", ps.ByName("vedioName"))
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)

	if err := json.Unmarshal(body, &vedioInfo); err == nil {
		fmt.Println("json.Unmarshal vedioInfo")
		fmt.Println(vedioInfo)
	} else {
		fmt.Println("json.Unmarshal err")
		fmt.Println(err)
	}
}

func deleteVedio(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var vedioInfo VedioInfo

	fmt.Fprintf(w, "deleteVedio, vedio name %s!\n", ps.ByName("vedioName"))
	fmt.Printf("deleteVedio, vedio name %s!\n", ps.ByName("vedioName"))
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)

	if err := json.Unmarshal(body, &vedioInfo); err == nil {
		fmt.Println("json.Unmarshal vedioInfo")
		fmt.Println(vedioInfo)
	} else {
		fmt.Println("json.Unmarshal err")
		fmt.Println(err)
	}
}

func purchaseVedio(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var vedioInfo VedioInfo

	fmt.Fprintf(w, "purchaseVedio, vedio name %s!\n", ps.ByName("vedioName"))
	fmt.Printf("purchaseVedio, vedio name %s!\n", ps.ByName("vedioName"))
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)

	if err := json.Unmarshal(body, &vedioInfo); err == nil {
		fmt.Println("json.Unmarshal vedioInfo")
		fmt.Println(vedioInfo)
	} else {
		fmt.Println("json.Unmarshal err")
		fmt.Println(err)
	}
}

func playVedio(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var vedioInfo VedioInfo

	fmt.Fprintf(w, "playVedio, vedio name %s!\n", ps.ByName("vedioName"))
	fmt.Printf("playVedio, vedio name %s!\n", ps.ByName("vedioName"))
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)

	if err := json.Unmarshal(body, &vedioInfo); err == nil {
		fmt.Println("json.Unmarshal vedioInfo")
		fmt.Println(vedioInfo)
	} else {
		fmt.Println("json.Unmarshal err")
		fmt.Println(err)
	}
}
