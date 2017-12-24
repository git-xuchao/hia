package config

import (
	"fmt"
	"testing"

	"github.com/Unknwon/goconfig"
)

func TestConfig(t *testing.T) {
	fmt.Println("TestConfig")
	var err error
	var cfg *goconfig.ConfigFile
	cfg, err = goconfig.LoadConfigFile("./conf.ini")
	if err != nil {
		fmt.Println("LoadConfigFile err", err)
		cfg, err = goconfig.LoadFromData([]byte(""))
		if err != nil {
			fmt.Println("LoadFromData err", err)
		}
	}

	cfg.SetValue("Demo", "key4", "hello girl!")
	goconfig.SaveConfigFile(cfg, "./conf.ini")
	/*
	 *v, err := cfg.GetValue("DataBase", "ip")
	 */
}
