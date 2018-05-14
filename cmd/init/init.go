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
	if os.Getuid() != 0 {
		panic("require root privilege")
	}

	syscall.Mount("proc", "/proc", "proc", 0, "")
	fmt.Printf("uid %d\n", os.Getuid())

	var cfg runc.Config
	config := os.Getenv("CONFIG")
	json.Unmarshal([]byte(config), &cfg)
	cfg.Rlimits = []runc.Rlimit{
		{
			Type: syscall.RLIMIT_CPU,
			Hard: 1,
			Soft: 1,
		},
	}
	err := setupRlimits(cfg.Rlimits)
	if err != nil {
		log.Printf("%s", err)
	}

	err = syscall.Exec("./main", []string{"./main"}, nil)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func setupRlimits(limits []runc.Rlimit) error {
	fmt.Printf("limits %#v\n", limits)
	for _, rlimit := range limits {
		if err := syscall.Setrlimit(rlimit.Type, &syscall.Rlimit{Max: rlimit.Hard, Cur: rlimit.Soft}); err != nil {
			return fmt.Errorf("fail to set rlimit %d: %s", rlimit.Type, err)
		}
	}
	return nil
}
