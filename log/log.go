package log

import (
	"fmt"
	log "github.com/inconshreveable/log15"
	"testing"
)

var Log = log.New()

func TestRedis2(t *testing.T) {
	handler := log.MultiHandler(
		log.Must.FileHandler("./app.log", log.LogfmtFormat()),
		/*
		 *log.StdoutHandler,
		 */
		log.StderrHandler)
	Log.SetHandler(handler)
	fmt.Println("test redis start")
	log.Info("Program startingasdfasdfasd")
}
