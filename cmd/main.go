package main

import (
	"log"
	"time"

	"github.com/sakeven/runc"
)

func main() {
	err := runc.C()
	if err != nil {
		log.Printf("%s", err)
		return
	}

	time.Sleep(time.Millisecond)
	log.Printf("succ")
}
