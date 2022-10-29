package lyr

import (
	"log"
	"os"
)

func IfErr(err error) {
	if err != nil {
		log.Fatal("ERR ", err)
		os.Exit(-1)
	}
}
