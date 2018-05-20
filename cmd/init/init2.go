package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"syscall"

	"github.com/sakeven/runc/pkg/runc"
)

func main() {
	runtime.LockOSThread()
	fmt.Printf("uid %d\n", os.Getuid())

	var cfg runc.Config
	config := os.Getenv("CONFIG")
	json.Unmarshal([]byte(config), &cfg)
	err := setupRlimits(cfg.Rlimits)
	if err != nil {
		log.Printf("%s", err)
	}

	args := os.Args
	cmd := args[1]

	err = syscall.Exec(cmd, args[1:], nil)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func setupRlimits(limits []runc.Rlimit) error {
	fmt.Printf("limits %#v\n", limits)
	var err error
	for _, rlimit := range limits {
		err = syscall.Setrlimit(rlimit.Type, &syscall.Rlimit{
			Max: rlimit.Hard,
			Cur: rlimit.Soft})
		if err != nil {
			return fmt.Errorf("fail to set rlimit %d: %s", rlimit.Type, err)
		}
	}
	return nil
}
