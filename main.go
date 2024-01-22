package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

type RunContext struct {
	Debug bool
}

type ExecCmd struct {
	Path string `arg:"" name:"path" help:"Path to the Go script" type:"path"`
}

func (e *ExecCmd) Run(runCtx *RunContext) error {
	i := interp.New(interp.Options{})

	err := i.Use(stdlib.Symbols)
	if err != nil {
		return err
	}

	_, err = i.EvalPath(e.Path)
	if err != nil {
		return err
	}

	v, err := i.Eval("syscheck.Execute")
	if err != nil {
		panic(err)
	}

	exec, ok := v.Interface().(func() ([]string, error))
	if !ok {
		panic(fmt.Errorf("conversion failed"))
	}

	r, err := exec()
	if err != nil {
		panic(err)
	}

	_ = r

	//_, err = i.Eval(`r == "darwin"`)
	//if err != nil {
	//	panic(err)
	//}

	return nil
}

var cli struct {
	Debug bool `help:"Enable debug mode."`

	Exec ExecCmd `cmd:"" help:"Execute Go script"`
}

func main() {
	runCtx := kong.Parse(&cli)

	// Call the Run() method of the selected parsed command.
	err := runCtx.Run(&RunContext{Debug: cli.Debug})
	runCtx.FatalIfErrorf(err)
}
