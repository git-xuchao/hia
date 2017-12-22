package ysdb

import (
	"fmt"
	"io/ioutil"
	"testing"

	"hia/cmd/spark/types"
)

func TestMySQLAddAuthorUser(t *testing.T) {

	var user types.User

	abi, _ := ioutil.ReadFile("./copyright_sol_copyright.abi")

	user.UserID = 987654325
	user.Password = "123456"
	/*
	 *user.UserType = "common"
	 */
	user.UserType = "author"
	user.UserName = "xuchao"
	user.UserIdCard = "xxxx"
	user.EthAccount = "0x6e2d604754ae054e2558b38a265cb84fccb975f6"
	user.EthContractAddr = "0xa231475d813a4e642c0f98fe3167211e2e9d133d"

	user.EthAbi = string(abi)

	// user.EthKey = "a"
	// user.Email = "a"

	db := NewDbMysql()

	err := db.UserAdd(&user)
	if err != nil {
		fmt.Println(err)
	}

	res, err1 := db.UserQuery(&user, "")

	fmt.Println(res, err1)

}
func TestMySQLAddCommonUser(t *testing.T) {

	var user types.User

	user.UserID = 987654326
	user.Password = "123456"
	user.UserType = "common"
	user.UserName = "alan"
	user.UserIdCard = "xxxx"
	user.EthAccount = "0x40b0eb2afe313d260b7574d4ffc6a130e9ad28ed"

	db := NewDbMysql()

	err := db.UserAdd(&user)
	if err != nil {
		fmt.Println(err)
	}

	res, err1 := db.UserQuery(&user, "")

	fmt.Println(res, err1)

}
