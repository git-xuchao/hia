package lab

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func TestFileInfo(t *testing.T) {
	myfolder := `/home/alan/tmp/data/node1/keystore`

	files, _ := ioutil.ReadDir(myfolder)
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			fmt.Println(file.Name())
			fmt.Println(strings.Contains(file.Name(), "47ed13c445022a9ecdf7577cdb11d186f8e5f8ce"))
			if strings.Contains(file.Name(), "47ed13c445022a9ecdf7577cdb11d186f8e5f8ce") == true {
				fmt.Println(myfolder + "/" + file.Name())
				dat, _ := ioutil.ReadFile(myfolder + "/" + file.Name())
				fmt.Print(string(dat))
			}
		}
	}
}
