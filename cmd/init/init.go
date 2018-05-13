package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"syscall"
)

func main() {
	runtime.LockOSThread()
	fmt.Printf("uid is %d\n", os.Getuid())

	syscall.Mount("proc", "/proc", "proc", 0, "")
	err := syscall.Exec("./main", []string{"./main"}, nil)
	if err != nil {
		log.Printf("%s", err)
	}
}
