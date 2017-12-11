package lab

import (
	"encoding/json"
	"fmt"
	"testing"
)

type User struct {
	UserName string `json:"username"`
	NickName string `json:"nickname"`
	Age      int
	Birthday string
	Sex      string
	Email    string
	Phone    string
}

/*结构体转json*/

func testStruct() (ret string, err error) {
	user1 := &User{
		UserName: "user1",
		NickName: "上课看似",
		Age:      18,
		Birthday: "2008/8/8",
		Sex:      "男",
		Email:    "mahuateng@qq.com",
		Phone:    "110",
	}

	data, err := json.Marshal(user1)
	if err != nil {
		fmt.Printf("json.marshal failed, err:", err)
		return
	}

	fmt.Printf("%s\n", string(data))
	ret = string(data)

	return
}

func TestJsonToStruct(t *testing.T) {
	data, err := testStruct()
	if err != nil {
		fmt.Println("test struct failed, ", err)
		return
	}

	var user1 User
	err = json.Unmarshal([]byte(data), &user1)
	if err != nil {
		fmt.Println("Unmarshal failed, ", err)
		return
	}
	fmt.Printf("TestJsonToStruct\n")
	fmt.Println(user1)
}

func TestStructToJson(t *testing.T) {
	testStruct()
	fmt.Println("----")
}
