package main

import (
	"log"
	"time"

	"github.com/sakeven/runc/consts"
	"github.com/sakeven/runc/pkg/compile"
	"github.com/sakeven/runc/pkg/runc"
)

var testCode = `
#include<stdio.h>
#include<stdlib.h>
#include <signal.h> 


void sigroutine(int dunno) {
  printf("quit %d\n", dunno);
  exit(1);
}

int main() {
  signal(SIGQUIT, sigroutine); 
  signal(SIGKILL, sigroutine); 
  signal(SIGTERM, sigroutine); 
  printf("hello world\n");
  printf("uid %d pid %d\n", getuid(), getpid());
  int i = 0;
  while (1) {
    i ++;
  }
  printf("%d", i);
  int *mem = (int *)malloc(sizeof(int)*1000*1000);
  if( mem == NULL) {
    printf("alloc failed");
    fflush(stdout);
    return 1;
  }
  mem[0] = 1;
  printf("alloc %d\n", mem[0]);
  fflush(stdout);
  return 0;
}
`

const STD_MB = 1048576
const STD_T_LIM = 2
const STD_F_LIM = (STD_MB << 5)
const STD_M_LIM = (STD_MB << 7)

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
		InitProcess: []string{"./init", "./main"},
		Rlimits: []runc.Rlimit{
			{
				Type: consts.RLIMIT_CPU,
				Hard: 1,
				Soft: 1,
			},
			{
				Type: consts.RLIMIT_FSIZE,
				Hard: consts.STD_F_LIM + consts.STD_MB,
				Soft: consts.STD_F_LIM,
			},
			{
				Type: consts.RLIMIT_NPROC,
				Hard: 1,
				Soft: 1,
			},
			{
				Type: consts.RLIMIT_STACK,
				Hard: consts.STD_MB << 6,
				Soft: consts.STD_MB << 6,
			},
			{
				Type: consts.RLIMIT_AS,
				Hard: consts.STD_MB << 4,
				Soft: consts.STD_MB << 4,
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
