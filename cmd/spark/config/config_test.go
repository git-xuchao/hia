package config

import (
	"fmt"
	"testing"

	"github.com/Unknwon/goconfig"
)

func TestConfig(t *testing.T) {
	fmt.Println("TestConfig")
	cfg, err := goconfig.LoadConfigFile("./conf.ini")
	if err != nil {
		fmt.Println("LoadConfigFile err")
		return
	}
	v, err := cfg.GetValue("DataBase", "ip")
	fmt.Printf("%s\n", v)
}
