package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/sakeven/runc/pkg/runc"
)

func main() {
	runtime.LockOSThread()
	if os.Getuid() != 0 {
		panic("require root privilege")
	}

	syscall.Mount("proc", "/proc", "proc", 0, "")
	runner := NewRunner()
	err := runner.Run()
	if err != nil {
		log.Printf("run init2 failed %s", err)
	}
}

type Runner struct {
	cfg *runc.Config
	cmd *exec.Cmd
}

func NewRunner() *Runner {
	runner := &Runner{
		cfg: getConfig(),
	}
	runner.cmd = newCommand(runner.cfg)
	return runner
}

func getConfig() *runc.Config {
	var cfg runc.Config
	config := os.Getenv("CONFIG")
	json.Unmarshal([]byte(config), &cfg)
	return &cfg
}

func newCommand(cfg *runc.Config) *exec.Cmd {
	args := os.Args[1:]
	cmd := exec.Command("./init2", args...)
	cmd.ExtraFiles = []*os.File{}
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	b, _ := json.Marshal(cfg)
	cmd.Env = []string{fmt.Sprintf("CONFIG=%s", b)}
	return cmd
}

// Run run.
func (rc *Runner) Run() error {
	t, m, err := rc.run()
	log.Printf("time %s mem %d, err %v", t, m, err)
	return err
}

func (rc *Runner) run() (time.Duration, int, error) {
	err := rc.cmd.Start()
	if err != nil {
		log.Printf("error %s", err)
		return 0, 0, err
	}

	topMem, _ := getMem(rc.cmd.Process.Pid)
	go func() {
		for {
			mem, err := getMem(rc.cmd.Process.Pid)
			if err != nil {
				log.Printf("get mem %s", err)
				return
			}
			if mem > topMem {
				topMem = mem
			}
		}
	}()

	time.Sleep(time.Microsecond)
	err = rc.cmd.Wait()

	ps := rc.cmd.ProcessState
	t := ps.SystemTime() + ps.UserTime()
	status := rc.cmd.ProcessState.Sys().(syscall.WaitStatus)
	log.Printf("exit status %d", status.Signal())
	log.Printf("exit status %d", ((uint32(status) & 0xff00) >> 8))
	return t, topMem, err
}

func getMem(pid int) (int, error) {
	f := fmt.Sprintf("/proc/%d/status", pid)
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return 0, err
	}

	var size int
	lines := strings.Split(string(b), "\n")
	for _, l := range lines {
		if strings.HasPrefix(l, "VmPeak") {
			l = strings.TrimPrefix(l, "VmPeak:")
			_, err = fmt.Sscanf(l, "%d", &size)
			return size, err
		}
	}
	return 0, fmt.Errorf("not found")
}
