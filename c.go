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

// C c.
func C() error {
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()

	// cmd := exec.CommandContext(ctx, "./cmain")

	cmd := exec.Command("./cmain")
	cmd.Dir = "./root"
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Chroot: "./root",
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
	}

	cmd.SysProcAttr.UidMappings = []syscall.SysProcIDMap{
		{ContainerID: 0, HostID: os.Getuid(), Size: 1},
	}
	cmd.SysProcAttr.GidMappings = []syscall.SysProcIDMap{
		{ContainerID: 0, HostID: os.Getgid(), Size: 1},
	}
	cmd.ExtraFiles = []*os.File{}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		log.Printf("error %s", err)
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	mem := getMem(cmd.Process.Pid)
	log.Printf("pid %d mem %s", cmd.Process.Pid, mem)
	go func() {
		for {
			mem := getMem(cmd.Process.Pid)
			if mem == "" {
				return
			}
			log.Printf("mem %s", mem)
		}
	}()
	time.Sleep(time.Microsecond)
	cmd.Wait()
	ps := cmd.ProcessState
	t := ps.SystemTime() + ps.UserTime()
	log.Printf("time %s", t)
	log.Printf("%#v", ps.SysUsage())
	return nil
}

func getMem(pid int) string {
	f := fmt.Sprintf("/proc/%d/status", pid)
	b, err := ioutil.ReadFile(f)
	if err != nil {
		log.Printf("read proc error %s", err)
		return ""
	}
	lines := strings.Split(string(b), "\n")
	for _, l := range lines {
		if strings.HasPrefix(l, "VmPeak") {
			return l
		}
	}
	return ""
}
