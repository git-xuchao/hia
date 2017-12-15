package myethclient

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

type EthClient struct {
	client  *rpc.Client
	abiFile *os.File
	abi     abi.ABI
	decoder *json.Decoder
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

func (this *EthClient) SetCallMsg(msg *ethereum.CallMsg, from string, to string, data []byte) error {
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

	msg.Data = data

	return nil
}

func (this *EthClient) CallContract(msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	var hex hexutil.Bytes
	err := this.client.CallContext(context.Background(), &hex, "eth_call", toCallArg(msg), toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}
	return hex, nil
}

func (this *EthClient) Call(result interface{}, method string, args ...interface{}) error {
	return this.client.Call(result, method, args...)
}

func (this *EthClient) Close() {
	this.client.Close()
}

func (this *EthClient) NewAccount(password string) string {
	var result json.RawMessage
	var account string

	this.Call(&result, "personal_newAccount", password)
	account = common.ToHex(result)

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
