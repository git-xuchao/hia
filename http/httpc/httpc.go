package httpc

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func TestHttpcCommand() {
	response, _ := http.Get("http://localhost:80/hello")
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
}
