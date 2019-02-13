package helpers

import (
	"fmt"
	"log"
	"os"
)

var debug *log.Logger

func init() {
	path, _ := os.LookupEnv("CPI_TRACE")
	if path == "" {
		return
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	debug = log.New(file, "", 0)
}

// Log ...
func Log(format string, values ...interface{}) {
	debug.Println(fmt.Sprintf(format, values...))
}

// Msg ...
func Msg(fmt string, values ...interface{}) {
	Log("... "+fmt, values...)
}
