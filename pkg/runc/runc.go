package runc

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// config of process
type Config struct {
	RootfsDir   string
	Chroot      bool
	InitProcess []string
	Rlimits     []Rlimit
}

// Rlimit represents resource limit
type Rlimit struct {
	Type int
	Hard uint64
	Soft uint64
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
		//		AmbientCaps: []uintptr{24},
	}
	if cfg.Chroot {
		cmd.SysProcAttr.Chroot = cfg.RootfsDir
	}
	cmd.ExtraFiles = []*os.File{}
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	b, _ := json.Marshal(cfg)
	cmd.Env = []string{fmt.Sprintf("CONFIG=%s", b)}
	return cmd
}

// Run run.
func (rc *Runc) Run() error {
	return rc.cmd.Start()
}
