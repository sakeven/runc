package runc

import "os"

type Runcer interface {
	Run() error
	//run() (time.Duration, int, error)
}

type Checker interface {
	Check(*os.File, *os.File) (int, error)
}

type Compiler interface {
	Compile(code string, lang int, dist string) error
}
