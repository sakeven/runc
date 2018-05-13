package main

import (
	"log"
	"time"

	"github.com/sakeven/runc/pkg/compile"
	"github.com/sakeven/runc/pkg/runc"
)

var testCode = `
#include<stdio.h>

int main() {
  printf("hello world\n");
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
