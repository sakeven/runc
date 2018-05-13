package compile

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type Compiler struct {
}

const (
	LangC   = 0
	LangCXX = 1
)

func (c *Compiler) Compile(code string, lang int, dist string) error {
	switch lang {
	case LangC:
		return CCompiler(code, dist)
	}
	return fmt.Errorf("not support lang %d", lang)
}

func CCompiler(code string, dist string) error {
	src := "./Main.c"
	err := prepareSoureCode(code, src)
	if err != nil {
		return err
	}
	return cCompile(src, dist)
}

func prepareSoureCode(code string, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(code)
	return err
}

func cCompile(src, dist string) error {
	// TODO timeout
	var b bytes.Buffer
	cmd := exec.Command("gcc", src, "-o", dist, "-Wall", "-lm", "--static", "-std=c99", "-DONLINE_JUDGE")
	cmd.Stderr = &b
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf(b.String())
	}
	return nil
}
