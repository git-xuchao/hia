// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package myethclient

import (
	"fmt"
	/*
	 *"reflect"
	 */
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
)

/*
 *func TestClientHTTP(t *testing.T) {
 *    server := newTestServer("service", new(rpc.Service))
 *    defer server.Stop()
 *    client := rpc.DialInProc(server)
 *    defer client.Close()
 *
 *    var resp rpc.Result
 *    if err := client.Call(&resp, "service_echo", "hello", 10, &rpc.Args{"world"}); err != nil {
 *        t.Fatal(err)
 *    }
 *
 *    fmt.Println(resp.Int, resp.String, resp.Args)
 *    if !reflect.DeepEqual(resp, rpc.Result{"hello", 10, &rpc.Args{"world"}}) {
 *        t.Errorf("incorrect result %#v", resp)
 *    }
 *}
 *
 *func newTestServer(serviceName string, service interface{}) *rpc.Server {
 *    server := rpc.NewServer()
 *    if err := server.RegisterName(serviceName, service); err != nil {
 *        panic(err)
 *    }
 *    return server
 *}
 */

type Result struct {
	String string
}

func TestNodeInfo(t *testing.T) {
	var result p2p.NodeInfo

	fmt.Println("TestNodeInfo")
	client, _ := rpc.Dial("http://127.0.0.1:8001")
	defer client.Close()

	if err := client.Call(&result, "admin_nodeInfo"); err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
	fmt.Println("\n")

	data, err := json.Marshal(result)
	if err != nil {
		fmt.Printf("json.marshal failed, err:", err)
		return
	}

	fmt.Printf("%s\n", string(data))
	fmt.Printf("enode:%s\n", result.Enode)
}

func TestClientShowNodeInfo(t *testing.T) {
	var client *EthClient
	client = NewEthClient()
	client.Dial("http://127.0.0.1:8001")
	defer client.Close()

	var result p2p.NodeInfo
	if err := client.Call(&result, "admin_nodeInfo"); err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
	fmt.Println("\n")

	data, err := json.Marshal(result)
	if err != nil {
		fmt.Printf("json.marshal failed, err:", err)
		return
	}

	fmt.Printf("%s\n", string(data))
	fmt.Printf("enode:%s\n", result.Enode)
}

func TestCallContact(t *testing.T) {
	var result json.RawMessage
	var mm ethereum.CallMsg
	var addr common.Address
	var v abi.ABI

	fmt.Println("TestCallContact")
	fp, _ := os.Open("./Hello_sol_Hello.abi")
	dec := json.NewDecoder(fp)

	for {
		err := dec.Decode(&v)
		if err != nil {
			break
		}
		fmt.Printf("%v\n", v)
	}

	ff, _ := v.Pack("getAuthorString", [32]byte{123}, false)
	fmt.Printf("data: 0x%x\n", ff)

	client, _ := ethclient.Dial("http" + "://" + "192.168.31.52:8545")
	temp := common.FromHex("0x192fc81ea2f59af885f2c55cf262cd77ec155335")
	slice := addr[:0]
	slice = append(slice, temp...)

	fmt.Printf("slice: %x\n", slice)

	mm.Data = ff
	mm.To = &addr

	fmt.Print("msg", mm)

	result, _ = client.CallContract(context.Background(), mm, nil)
	fmt.Printf("result %v\n", result)
}

func TestCallContact2(t *testing.T) {
	var cli *EthClient
	var msg ethereum.CallMsg

	cli = NewEthClient()
	cli.Dial("http" + "://" + "192.168.31.52:8545")
	cli.ConstructAbi("./Hello_sol_Hello.abi")
	data, _ := cli.PackMethod("getAuthorString", [32]byte{123}, false)
	fmt.Println("data:", data)
	cli.SetCallMsg(&msg, "", "0x192fc81ea2f59af885f2c55cf262cd77ec155335", "", "", "", data)
	fmt.Println("msg", msg)
	result, _ := cli.CallContract2(msg, nil)
	fmt.Printf("result %v\n", result)
}

func TestSendTransaction(t *testing.T) {
	var cli *EthClient
	var msg ethereum.CallMsg

	fmt.Println("TestSendTransaction")
	cli = NewEthClient()
	cli.Dial("http" + "://" + "192.168.31.52:8545")
	cli.ConstructAbi("./Hello_sol_Hello.abi")
	data, _ := cli.PackMethod("getAuthorString", [32]byte{123}, false)
	fmt.Println("data:", data)
	cli.SetCallMsg(&msg, "0x6d83edbcc8c55c183a0695bb39ac20e7cf17f100", "0x192fc81ea2f59af885f2c55cf262cd77ec155335", "", "", "", data)
	fmt.Println("msg", msg)
	result, err := cli.SendTransaction(msg, "123456")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("result %v\n", result)
	}
}

func TestCallContactMethod(t *testing.T) {
	var cli *EthClient
	var msg ethereum.CallMsg

	fmt.Println("TestCallContactMethod")
	cli = NewEthClient()
	cli.Dial("http" + "://" + "192.168.31.52:8545")
	cli.ConstructAbi("./Hello_sol_Hello.abi")

	cli.SetCallMsg(&msg, "0x6d83edbcc8c55c183a0695bb39ac20e7cf17f100", "0x192fc81ea2f59af885f2c55cf262cd77ec155335", "", "", "", nil)

	result, err := cli.CallContractMethod(msg, "123456", "getAuthorString", [32]byte{123}, false)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("result %v\n", result)
	}
}

func TestCallContactMethod2(t *testing.T) {
	var cli *EthClient
	var msg ethereum.CallMsg

	fmt.Println("TestCallContactMethod2")
	cli = NewEthClient()
	cli.Dial("http" + "://" + "192.168.31.52:8545")
	cli.ConstructAbi("./copyright_sol_copyright.abi")

	cli.SetCallMsg(&msg, "0x6e2d604754ae054e2558b38a265cb84fccb975f6", "0xa231475d813a4e642c0f98fe3167211e2e9d133d", "", "", "", nil)

	result, err := cli.CallContractMethod(msg, "123456", "purchaseVideo", "xuchao", "dsfasdfdsf")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("result %v\n", result)
	}
}

func TestCallContractMethodOnly(t *testing.T) {
	var cli *EthClient
	var msg ethereum.CallMsg

	cli = NewEthClient()
	cli.Dial("http" + "://" + "192.168.31.52:8545")
	cli.ConstructAbi("./copyright_sol_copyright.abi")
	data, _ := cli.PackMethod("playVideo", "xuchao2", "dsfasdfasdasdfdfddsfasdfsf")
	fmt.Println("data:", data)
	cli.SetCallMsg(&msg, "", "0xa231475d813a4e642c0f98fe3167211e2e9d133d", "", "", "", data)
	fmt.Println("msg", msg)
	result, _ := cli.CallContract2(msg, nil)
	fmt.Printf("result %v\n", result)
}

func TestCallContractMethodOnly2(t *testing.T) {
	var cli *EthClient
	var msg ethereum.CallMsg

	fmt.Println("CallContractMethodOnly2")
	cli = NewEthClient()
	cli.Dial("http" + "://" + "192.168.31.52:8545")
	cli.ConstructAbi("./copyright_sol_copyright.abi")

	cli.SetCallMsg(&msg, "0x6e2d604754ae054e2558b38a265cb84fccb975f6", "0xa231475d813a4e642c0f98fe3167211e2e9d133d", "", "", "", nil)

	result, err := cli.CallContractMethodOnly(msg, nil, "playVideo", "xuchao", "dsfadsf")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("result %v\n", result)
	}
}

func TestCallContractMethodPack(t *testing.T) {
	var cli *EthClient
	var msg ethereum.CallMsg

	fmt.Println("TestCallContractMethodPack")
	cli = NewEthClient()
	cli.Dial("http" + "://" + "192.168.31.52:8545")
	cli.ConstructAbi("./copyright_sol_copyright.abi")

	cli.SetCallMsg(&msg, "0x6e2d604754ae054e2558b38a265cb84fccb975f6", "0xa231475d813a4e642c0f98fe3167211e2e9d133d", "", "", "", nil)

	/*
	 *result, err := cli.CallContractMethodOnly(msg, nil, "playVideo", "alan", "abc.flv")
	 */
	result, err := cli.CallContractMethodPack(msg, "123456", "purchaseVideo", "alan", "abc.flv")
	/*
	 *result, err := cli.CallContractMethodPack(msg, "123456", "uploadVideo", "abc.flv")
	 */
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("result %v\n", result)
	}
}

func TestClientNewAccount(t *testing.T) {
	var client *EthClient
	client = NewEthClient()
	client.Dial("http://127.0.0.1:8001")
	defer client.Close()

	client.NewAccount("123456")
}

func TestClientListAccounts(t *testing.T) {
	var client *EthClient
	client = NewEthClient()
	client.Dial("http://127.0.0.1:8001")
	defer client.Close()

	client.ListAccounts()
}
