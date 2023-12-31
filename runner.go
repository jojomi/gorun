package gorun

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Runner struct {
	command       Command
	workingDir    string
	env           map[string]string
	nonZeroExitOK bool
	doLog         bool
	stdout        bool
	stderr        bool
}

func New() *Runner {
	return (&Runner{
		env: make(map[string]string, 0),
	}).Reset()
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

func (x *Runner) InWorkingDir(workingDir string) *Runner {
	x.workingDir = workingDir
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

func (x *Runner) AddEnv(key, value string) *Runner {
	x.env[key] = value
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

	if x.workingDir != "" {
		cmd.Dir = x.workingDir
	}
	if len(x.env) > 0 {
		cmd.Env = append(os.Environ(), x.makeEnv(x.env)...)
	}

	if !x.stdout {
		cmd.Stdout = io.MultiWriter(rr.stdoutBuffer, rr.combinedBuffer)
	} else {
		cmd.Stdout = io.MultiWriter(os.Stdout, rr.stdoutBuffer, rr.combinedBuffer)
	}

	if !x.stderr {
		cmd.Stderr = io.MultiWriter(rr.stderrBuffer, rr.combinedBuffer)
	} else {
		cmd.Stderr = io.MultiWriter(os.Stderr, rr.stderrBuffer, rr.combinedBuffer)
	}

	cmd.Stdin = os.Stdin

	// logging
	if x.doLog {
		fmt.Println(x.command.String())
		fmt.Println("ENV: " + strings.Join(x.makeEnv(x.env), ", "))
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

func (x *Runner) makeEnv(env map[string]string) []string {
	result := make([]string, 0, len(env))
	for k, v := range env {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}
	return result
}
