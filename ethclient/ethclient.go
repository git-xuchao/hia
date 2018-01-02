package myethclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

type EthClient struct {
	client                *rpc.Client
	abiFile               *os.File
	abi                   abi.ABI
	decoder               *json.Decoder
	keyStoreSearchingPath string
}

func (this *EthClient) Dial(addr string) error {
	var err error

	c, err := rpc.Dial(addr)
	if err != nil {
		return err
	}
	this.client = c

	return err
}

/*
 *func (this *EthClient) OpenAbiFile(path string) error {
 *    fp, err := os.Open(path)
 *    this.abiFile = fp
 *
 *    return err
 *}
 *
 *func (this *EthClient) NewDecoder() error {
 *    this.decoder = json.NewDecoder(this.abiFile)
 *
 *    return nil
 *}
 *
 *func (this *EthClient) DecodeAbi() error {
 *    var err error
 *    for {
 *        err = this.decoder.Decode(&this.abi)
 *        if err != nil {
 *            break
 *        }
 *    }
 *    fmt.Printf("%v\n", this.abi)
 *
 *    return err
 *}
 */

func (this *EthClient) ConstructAbi2(def string) error {
	var err error

	this.abi, err = abi.JSON(strings.NewReader(def))
	fmt.Printf("%v\n", this.abi)

	return err
}

func (this *EthClient) ConstructAbi(path string) error {
	fp, err := os.Open(path)
	this.abiFile = fp
	this.decoder = json.NewDecoder(fp)
	for {
		err = this.decoder.Decode(&this.abi)
		if err != nil {
			break
		}
	}
	fmt.Printf("%v\n", this.abi)

	return err
}

func (this *EthClient) PackMethod(name string, args ...interface{}) ([]byte, error) {
	return this.abi.Pack(name, args...)
}

func (this *EthClient) SetCallMsg(msg *ethereum.CallMsg, from string, to string, gas string, gasPrice string, value string, data []byte) error {
	var addr common.Address

	if len(from) != 0 {
		slice := msg.From[:0]
		temp := common.FromHex(from)
		slice = append(slice, temp...)
	}

	if len(to) != 0 {
		slice := addr[:0]
		temp := common.FromHex(to)
		slice = append(slice, temp...)
		msg.To = &addr
	}

	if gas != "" {
		msg.Gas = new(big.Int)
		msg.Gas.SetString(gas, 0)
	}

	if gasPrice != "" {
		msg.GasPrice = new(big.Int)
		msg.GasPrice.SetString(gasPrice, 0)
	}

	if value != "" {
		msg.Value = new(big.Int)
		msg.Value.SetString(value, 0)
	}

	msg.Data = data

	return nil
}

func (this *EthClient) CallContract2(msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	var hex hexutil.Bytes
	err := this.client.CallContext(context.Background(), &hex, "eth_call", toCallArg(msg), toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}
	return hex, nil
}

func (this *EthClient) CallContractMethodOnly(msg ethereum.CallMsg, blockNumber *big.Int, method string, args ...interface{}) ([]byte, error) {
	var hex hexutil.Bytes
	var Err uint8
	var String string

	data, _ := this.PackMethod(method, args...)
	msg.Data = data

	err := this.client.CallContext(context.Background(), &hex, "eth_call", toCallArg(msg), toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}

	result := []interface{}{&Err, &String}

	err = this.abi.Unpack(&result, method, hex)
	if err != nil {
		return nil, err
	}

	if Err != 0 {
		return nil, errors.New(String)
	}

	return nil, err
}

func (this *EthClient) SendTransaction(msg ethereum.CallMsg, password string) ([]byte, error) {
	var hex hexutil.Bytes
	err := this.client.CallContext(context.Background(), &hex, "personal_sendTransaction", toCallArg(msg), password)
	if err != nil {
		fmt.Printf("error %v\n", err)
		return nil, err
	}
	return hex, nil
}

func (this *EthClient) GetTransactionReceipt(txHash string) (types.Receipt, error) {
	var hex types.Receipt
	err := this.client.CallContext(context.Background(), &hex, "eth_getTransactionReceipt", txHash)
	if err != nil {
		fmt.Printf("error %v\n", err)
		return hex, err
	}
	return hex, nil
}

func (this *EthClient) CallContractMethod(msg ethereum.CallMsg, password string, method string, args ...interface{}) ([]byte, error) {
	data, _ := this.PackMethod(method, args...)
	msg.Data = data
	return this.SendTransaction(msg, password)
}

func (this *EthClient) CallContractMethodPack(msg ethereum.CallMsg, password string, method string, args ...interface{}) ([]byte, error) {
	_, err := this.CallContractMethodOnly(msg, nil, method, args...)
	if err != nil {
		return nil, err
	} else {
		return this.CallContractMethod(msg, password, method, args...)
	}
}

func (this *EthClient) Call(result interface{}, method string, args ...interface{}) error {
	return this.client.Call(result, method, args...)
}

func (this *EthClient) SetKeyStoreSearchingPath(path string) error {
	this.keyStoreSearchingPath = path

	return nil
}

func (this *EthClient) GetKeyFileName(account string) (string, error) {
	fmt.Println("keypath\n", this.keyStoreSearchingPath)
	files, _ := ioutil.ReadDir(this.keyStoreSearchingPath)
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			fmt.Println(file.Name(), account)
			if strings.Contains(file.Name(), account) == true {
				fmt.Println("found key file name\n", file.Name())
				return file.Name(), nil
			}
		}
	}
	return "", nil
}

func (this *EthClient) GetKey(account string) (string, error) {
	fmt.Printf("EthClient GetKey, account:%s\n", account)
	fileName, _ := this.GetKeyFileName(account)
	fmt.Printf("Key File Name:%s\n", fileName)
	if fileName != "" {
		dat, _ := ioutil.ReadFile(this.keyStoreSearchingPath + "/" + fileName)
		return string(dat), nil
	} else {
		return "", nil
	}
}

func (this *EthClient) NewAccount(password string) string {
	var result json.RawMessage
	var account string

	this.Call(&result, "personal_newAccount", password)
	/*
	 *account = common.ToHex(result)
	 *fmt.Printf("new account:%s\n", account)
	 */
	account = string(result)

	fmt.Printf("new account:%s\n", account)

	return account
}

func (this *EthClient) ListAccounts() []string {
	/*
	 *var accounts []string
	 */
	accounts := make([]string, 4)

	this.Call(&accounts, "personal_listAccounts")
	/*
	 *fmt.Println(accounts)
	 */

	for index, value := range accounts {
		fmt.Printf("index=%d, value=%s\n", index, value)
	}

	fmt.Printf("slice len=%d, cap=%d", len(accounts), cap(accounts))

	return accounts
}

func (this *EthClient) StartMining() error {
	err := this.Call(nil, "miner_start", 4)
	if err != nil {
		fmt.Printf("error %v\n", err)
		return err
	}

	return err
}

func (this *EthClient) StopMining() error {
	err := this.Call(nil, "miner_stop")
	if err != nil {
		fmt.Printf("error %v\n", err)
		return err
	}

	return err
}

func (this *EthClient) IsMining() (bool, error) {
	var result bool

	err := this.Call(&result, "eth_mining")
	if err != nil {
		fmt.Printf("error %v\n", err)
		return result, err
	}

	return result, err
}

func (this *EthClient) Close() {
	this.client.Close()
}

func NewEthClient() *EthClient {
	return &EthClient{}
}

func toCallArg(msg ethereum.CallMsg) interface{} {
	arg := map[string]interface{}{
		"from": msg.From,
		"to":   msg.To,
	}
	if len(msg.Data) > 0 {
		arg["data"] = hexutil.Bytes(msg.Data)
	}
	if msg.Value != nil {
		arg["value"] = (*hexutil.Big)(msg.Value)
	}
	if msg.Gas != nil {
		arg["gas"] = (*hexutil.Big)(msg.Gas)
	}
	if msg.GasPrice != nil {
		arg["gasPrice"] = (*hexutil.Big)(msg.GasPrice)
	}
	return arg
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	return hexutil.EncodeBig(number)
}
