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
	if os.Getuid() != 0 {
		panic("require root privilege")
	}

	syscall.Mount("proc", "/proc", "proc", 0, "")

	// TODO set rlimit
	setupRlimits(nil)
	// TODO drop privilege

	err := syscall.Exec("./main", []string{"./main"}, nil)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
	os.Exit(0)
}

type Rlimit struct {
	Type int
	Hard uint64
	Soft uint64
}

func setupRlimits(limits []Rlimit) error {
	for _, rlimit := range limits {
		if err := syscall.Setrlimit(rlimit.Type, &syscall.Rlimit{Max: rlimit.Hard, Cur: rlimit.Soft}); err != nil {
			return fmt.Errorf("fail to set rlimit %d: %s", rlimit.Type, err)
		}
	}
	return nil
}
