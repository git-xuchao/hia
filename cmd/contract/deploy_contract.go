package contract

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/urfave/cli.v1"
	"hia/cmd/spark/types"
	"hia/cmd/spark/ysdb"
	"hia/ethclient"
)

// TODO: temp func
func getKeyAndFileName(account string) (string, string, error) {
	keyStoreSearchingPath := "/home/seagull/projects/eth_chain1/node1/keystore"
	var filename string
	files, _ := ioutil.ReadDir(keyStoreSearchingPath)
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			if strings.Contains(file.Name(), account) == true {
				filename = file.Name()
				break
			}
		}
	}
	if filename != "" {
		dat, _ := ioutil.ReadFile(keyStoreSearchingPath + "/" + filename)
		return filename, string(dat), nil
	} else {
		return "", "", nil
	}
}

func execShell(cmd string) (string, error) {
	shellCmd := exec.Command("/bin/bash", "-c", cmd)
	var out bytes.Buffer
	shellCmd.Stdout = &out

	err := shellCmd.Run()

	return out.String(), err
}

func ExecCompileContract(solPath string) (*[]string, error) {
	var err error
	pathArr := make([]string, 2)

	cmd := "solcjs --optimize --combined-json --abi --bin " + solPath
	_, err = execShell(cmd)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	split := strings.Split(solPath, "/")
	filename := split[len(split)-1]
	prefix1 := make([]byte, 1024)
	out, err := execShell("pwd")
	for i := 0; i < len(out); i++ {
		if out[i] == '\n' {
			break
		}
		prefix1[i] = out[i]
	}
	prefix := string(prefix1) + "/"
	binFilename := strings.Replace(filename, "sol", "bin", -1)
	abiFilename := strings.Replace(filename, "sol", "abi", -1)
	solFilePath := strings.Replace(solPath, "/", "_", -1)
	solFilePath = strings.Replace(solFilePath, ".", "_", -1)
	pathArr[0] = prefix + solFilePath + "_" + abiFilename
	pathArr[1] = prefix + solFilePath + "_" + binFilename
	fmt.Printf("abi[%s] bin[%s]\n", pathArr[0], pathArr[1])

	return &pathArr, err
}

func getUserInfoFromDB(db ysdb.YsDb, ethcli *myethclient.EthClient, user *types.User) error {
	var err error
	var resUser types.User

	resUser, err = db.UserQuerySimple(user)
	if err != nil {
		fmt.Printf("err %v\n", err)
		return err

		user.Password = "123456"
		account := ethcli.NewAccount(user.Password)
		len := len(account)
		user.EthKeyFileName, user.EthKey, _ = getKeyAndFileName(account[3 : len-1])
		user.UserType = "author"
		user.EthAccount = account
		fmt.Printf("new user info: %v\n", user)

		err = db.UserAdd(user)
		if err != nil {
			fmt.Printf("database add user error\n")
			return err
		}
	} else {
		user.EthAccount = resUser.EthAccount
		user.Password = resUser.Password
	}

	return nil
}

func DeployContract(userID string, abiPath string, binPath string) error {
	var err error

	// dial to eth console
	var ethcli *myethclient.EthClient
	ethcli = myethclient.NewEthClient()
	ethcli.Dial("http" + "://" + "192.168.31.30:8000")
	defer ethcli.Close()

	// connect db
	var db ysdb.YsDb
	db = ysdb.NewDbMysql()
	db.Init("mysql", "root:root@tcp(192.168.31.19)/test")

	// get user info from db
	var user types.User
	user.UserID, _ = strconv.ParseUint(userID, 10, 64)
	user.UserName = userID // TODO:userID cannot duplicate
	fmt.Print("before getUserInfoFromDB\n")
	err = getUserInfoFromDB(db, ethcli, &user)
	fmt.Print("after getUserInfoFromDB\n")
	if err != nil {
		return err
	}

	// send transaction
	var msg ethereum.CallMsg
	gas := "10000000" // TODO: set max in further test

	cmd1 := "cat " + binPath
	cmd2 := "cat " + abiPath
	binCode, err := execShell(cmd1)
	if err != nil {
		return err
	}
	abiCode, err := execShell(cmd2)
	if err != nil {
		return err
	}

	byteData := common.FromHex(binCode)
	len := len(user.EthAccount)
	ethAccount := string(user.EthAccount[3 : len-1])
	ethcli.SetCallMsg(&msg, ethAccount, "", gas, "0", "", byteData)
	transResult, err := ethcli.SendTransaction(msg, user.Password)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	} else {
		fmt.Printf("result %v\n", transResult)
	}

	// wait for contract addr
	isMining, err := ethcli.IsMining()
	if isMining == false {
		err = ethcli.StartMining()
		if err != nil {
			fmt.Printf("StartMining err: %v\n", err)
		} else {
			fmt.Printf("start mining\n")
		}
	}

	transAddr := common.ToHex(transResult)
	for {
		fmt.Printf("begin to get transaction receipt\n")
		receipt, err := ethcli.GetTransactionReceipt(transAddr)

		if err != nil {
			fmt.Printf("err: %v\n", err)
		} else {
			if isMining == false {
				err = ethcli.StopMining()
				if err != nil {
					fmt.Printf("StopMining err: %v\n", err)
				} else {
					fmt.Printf("stop mining\n")
				}
			}

			// add contract addr and abi to db
			var contractAddr common.Address
			contractAddr = receipt.ContractAddress
			user.EthContractAddr = contractAddr.Hex()
			user.EthAbi = abiCode
			err = db.UserUpdate(&user)
			if err != nil {
				return err
			}
			resUser, _ := db.UserQuerySimple(&user)
			fmt.Printf("user info: %v\n", resUser)
			break
		}
		time.Sleep(1e+09)
	}

	return err
}

func DeployContractPipeline(userID string, solPath string) error {
	path, err := ExecCompileContract(solPath)
	if err != nil {
		return err
	} else if *path == nil || len(*path) != 2 {
		fmt.Print("file path params err\n")
		return nil
	}

	err = DeployContract(userID, (*path)[0], (*path)[1])
	if err != nil {
		return err
	}

	return nil
}

func RunDeployContractPipeline(c *cli.Context) error {
	if c.NArg() != 2 {
		fmt.Printf("Warning: %d args are given, but 2 args are needed\n", c.NArg())
		fmt.Println("Usage: ./main contract total userID contractAbsPath\n")
		return nil
	}

	err := DeployContractPipeline(c.Args().Get(0), c.Args().Get(1))

	return err
}

func RunDeployContractOnly(c *cli.Context) error {
	if c.NArg() != 3 {
		fmt.Printf("Warning: %d args are given, but 3 args are needed\n", c.NArg())
		fmt.Println("Usage: ./main contract deploy userID abiAbsPath binAbsPath\n")
		return nil
	}

	err := DeployContract(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2))

	return err
}
