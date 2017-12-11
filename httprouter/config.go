package myhttprouter

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
)

type Route struct {
	Name    string
	Method  string
	Pattern string
	Handle  httprouter.Handle
}

type Routes []Route

var routes = Routes{
	Route{"Index", "GET", "/", Index},
	Route{"Hello", "GET", "/hello/:name", Hello},
	Route{"registration", "POST", "/users/:usertype", register},
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func register(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "register, usertype %s!\n", ps.ByName("usertype"))
	fmt.Printf("register, usertype %s!\n", ps.ByName("usertype"))
	body, _ := ioutil.ReadAll(r.Body)
	//    r.Body.Close()
	body_str := string(body)
	fmt.Println(body_str)
	var user User

	if err := json.Unmarshal(body, &user); err == nil {
		fmt.Println(user)
		user.Age += 100
		fmt.Println(user)
		ret, _ := json.Marshal(user)
		fmt.Fprint(w, string(ret))
	} else {
		fmt.Println(err)
	}
}
