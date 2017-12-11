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

package test

import (
	"fmt"
	/*
	 *"reflect"
	 */
	"testing"

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

func TestClientHTTP(t *testing.T) {
	client, _ := rpc.Dial("http://127.0.0.1:8001")
	defer client.Close()

	var resp interface{}
	if err := client.Call(&resp, "admin_nodeInfo"); err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}
