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

	"github.com/sakeven/runc/consts"
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
	t, m, status, err := rc.run()
	log.Printf("time %s mem %d status %d, err %v", t, m, status, err)
	return err
}

func (rc *Runner) run() (time.Duration, int, int, error) {
	limit := struct {
		Time   time.Duration
		Memory int
	}{
		Time:   time.Second,
		Memory: 32768, //kB
	}
	err := rc.cmd.Start()
	if err != nil {
		log.Printf("error %s", err)
		return 0, 0, consts.JudgeNA, err
	}

	var wstatus syscall.WaitStatus
	var topTime time.Duration
	var topMem int
	for {
		m, err := getTime(rc.cmd.Process.Pid)
		if err != nil {
			log.Printf("get time %s", err)
			break
		}
		if m > topTime {
			topTime = m
		}

		if topTime > limit.Time {
			return topTime, topMem, consts.JudgeTLE, nil
		}

		mem, err := getMem(rc.cmd.Process.Pid)
		if err != nil {
			log.Printf("get mem %s", err)
			break
		}
		if mem > topMem {
			topMem = mem
		}
		if topMem >= limit.Memory {
			return topTime, topMem, consts.JudgeMLE, nil
		}
	}

	time.Sleep(time.Microsecond)
	err = rc.cmd.Wait()
	ps := rc.cmd.ProcessState
	t := ps.SystemTime() + ps.UserTime()
	wstatus = ps.Sys().(syscall.WaitStatus)
	signal := wstatus.Signal()

	log.Printf("signal is %s", signal)
	// if wait status signal is one of
	// SIGCHLD, SIGALRM, SIGKILL, SIGXCPU, it should be tle.
	// if wait status signal is SIGXFSZ, it should be ole.
	// if wait status signal is others, it shoulde be re.
	var status int
	switch signal {
	case syscall.SIGCHLD, syscall.SIGALRM, syscall.SIGKILL, syscall.SIGXCPU:
		status = consts.JudgeTLE
	case syscall.SIGXFSZ:
		status = consts.JudgeOLE
	case syscall.Signal(-1):
		// exit normally
		status = consts.JudgeAC
	default:
		status = consts.JudgeRE
		// TODO record exactly signal
	}
	return t, topMem, status, err
}

// KB, error
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

func getTime(pid int) (time.Duration, error) {
	f := fmt.Sprintf("/proc/%d/stat", pid)
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return 0, err
	}

	var (
		name, status                             string
		ppid, pgid, sid, tty_nr, tty_pgrp, flags int
		min_flt, cmin_flt, maj_flt, cmaj_flt     int
		utime, stime                             int
	)

	fmt.Sscanf(string(b), "%d %s %s %d %d %d %d %d %d %d %d %d %d %d %d",
		&pid, &name, &status, &ppid, &pgid, &sid,
		&tty_nr,
		&tty_pgrp,
		&flags,
		&min_flt,
		&cmin_flt,
		&maj_flt,
		&cmaj_flt,
		&utime,
		&stime,
	)
	// utime and stime is in jiffies unit.
	// jiffies = seconds * Hz
	return time.Duration(utime+stime) * time.Second / 100, nil
}
