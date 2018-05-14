package main

import (
	"log"
	"syscall"
	"time"

	"github.com/sakeven/runc/pkg/compile"
	"github.com/sakeven/runc/pkg/runc"
)

var testCode = `
#include<stdio.h>
#include<stdlib.h>

int main() {
  printf("hello world\n");
  printf("uid %d\n", getuid());
  sleep(1);
  int i = 0;
  while(1) {
    i++;
  }
  printf("%d\n", i);
  return 0;
}
`

func main() {
	compiler := &compile.Compiler{}
	err := compiler.Compile(testCode, compile.LangC, "./root/main")
	if err != nil {
		log.Printf("%s", err)
		return
	}

	cfg := &runc.Config{
		RootfsDir:   "./root",
		Chroot:      true,
		InitProcess: []string{"./init"},
		Rlimits: []runc.Rlimit{
			{
				Type: syscall.RLIMIT_CPU,
				Hard: 1,
				Soft: 1,
			},
		},
	}
	runner := runc.New(cfg)
	err = runner.Run()
	if err != nil {
		log.Printf("%s", err)
		return
	}

	// TODO output compare

	time.Sleep(time.Millisecond)
	log.Printf("succ")
}
