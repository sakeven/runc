package runc

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// config of process
type Config struct {
	RootfsDir   string
	Chroot      bool
	InitProcess []string
}

// Runc handles container process.
type Runc struct {
	cmd    *exec.Cmd
	config *Config
}

// New creates a new Runc instance.
func New(cfg *Config) *Runc {
	return &Runc{
		cmd: newCommand(cfg),
	}
}

func newCommand(cfg *Config) *exec.Cmd {
	cmd := exec.Command(cfg.InitProcess[0], (cfg.InitProcess[1:])...)
	cmd.Dir = cfg.RootfsDir
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getuid(), Size: 1},
		},
		GidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getgid(), Size: 1},
		},
	}
	if cfg.Chroot {
		cmd.SysProcAttr.Chroot = cfg.RootfsDir
	}
	cmd.ExtraFiles = []*os.File{}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

// Run run.
func (rc *Runc) Run() error {
	t, m, err := rc.run()
	log.Printf("time %s mem %d, err %v", t, m, err)
	return err
}

func (rc *Runc) run() (time.Duration, int, error) {
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
				log.Printf("%s", err)
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
