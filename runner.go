package gorun

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Runner struct {
	command       Command
	nonZeroExitOK bool
	doLog         bool
	stdout        bool
	stderr        bool
}

func New() *Runner {
	return (&Runner{}).Reset()
}

func NewLocal(localCommand string) *Runner {
	return NewWithCommand(LocalCommandFrom(localCommand))
}

func NewWithCommand(command Command) *Runner {
	return New().WithCommand(command)
}

func (x *Runner) Reset() *Runner {
	x.doLog = false
	x.nonZeroExitOK = false
	x.stdout = true
	x.stderr = true

	return x
}

func (x *Runner) WithCommand(command Command) *Runner {
	x.command = command
	return x
}

func (x *Runner) Silent() *Runner {
	x.WithoutStdout()
	x.WithoutStderr()
	return x
}

func (x *Runner) WithoutStdout() *Runner {
	x.stdout = false
	return x
}

func (x *Runner) WithoutStderr() *Runner {
	x.stderr = false
	return x
}

func (x *Runner) NonZeroExitOK() *Runner {
	x.nonZeroExitOK = true
	return x
}

func (x *Runner) LogCommand(value bool) *Runner {
	x.doLog = value
	return x
}

func (x *Runner) MustExec() *RunResult {
	result, err := x.Exec()
	if err != nil {
		panic(err)
	}
	return result
}

func (x *Runner) Exec() (*RunResult, error) {
	err := x.validate()
	if err != nil {
		return nil, err
	}

	rr := NewRunResult()
	rr.nonZeroExitOK = x.nonZeroExitOK

	// execute
	cmd := exec.Command(x.command.Binary(), x.command.Args()...)
	rr.Cmd = cmd

	/* TODO
	cmd.Dir = c.workingDir
	cmd.Env = c.GetFullEnv()
	*/

	if !x.stdout {
		cmd.Stdout = rr.stdoutBuffer
	} else {
		cmd.Stdout = io.MultiWriter(os.Stdout, rr.stdoutBuffer)
	}

	if !x.stderr {
		cmd.Stderr = rr.stderrBuffer
	} else {
		cmd.Stderr = io.MultiWriter(os.Stderr, rr.stderrBuffer)
	}

	cmd.Stdin = os.Stdin

	// logging
	if x.doLog {
		fmt.Println(x.command.String())
	}

	err = cmd.Start()
	if err != nil {
		rr.err = err
		return rr, nil
	}
	rr.Process = cmd.Process

	err = rr.Cmd.Wait()
	rr.ProcessState = rr.Cmd.ProcessState
	rr.ProcessError = err

	return rr, nil
}

func (x *Runner) ExecAsync() {
	// TODO implement ExecAsync
}

func (x *Runner) validate() error {
	if x.command == nil {
		return fmt.Errorf("missing command, use WithCommand")
	}
	return nil
}
