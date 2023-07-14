package gorun

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

type RunResult struct {
	command       Command
	nonZeroExitOK bool
	err           error
	Cmd           *exec.Cmd
	Process       *os.Process
	ProcessState  *os.ProcessState
	ProcessError  error
	stdoutBuffer  *bytes.Buffer
	stderrBuffer  *bytes.Buffer
}

func NewRunResult() *RunResult {
	p := &RunResult{}
	p.stdoutBuffer = bytes.NewBuffer(make([]byte, 0, 100))
	p.stderrBuffer = bytes.NewBuffer(make([]byte, 0, 100))
	return p
}

func (x RunResult) Successful() bool {
	if x.nonZeroExitOK {
		return true
	}
	ex, err := x.ExitCode()
	if err != nil {
		return false
	}
	return ex == 0
}

func (x RunResult) Failed() bool {
	return !x.Successful()
}

// CombinedOutput returns a string representation of all the output of the process denoted
// by this struct.
func (x RunResult) CombinedOutput() (string, error) {
	out, err := x.Cmd.CombinedOutput()
	return string(out), err
}

// MustCombinedOutput returns a string representation of all the output of the process denoted
// by this struct.
func (x RunResult) MustCombinedOutput() string {
	out, err := x.CombinedOutput()
	if err != nil {
		panic(err)
	}
	return out
}

// Stdout returns a string representation of the output of the process denoted
// by this struct.
func (x RunResult) Stdout() string {
	return x.stdoutBuffer.String()
}

// StdoutTrimmed returns a string representation of the output of the process denoted
// by this struct with surrounding whitespace removed.
func (x RunResult) StdoutTrimmed() string {
	return strings.TrimSpace(x.Stdout())
}

// StderrTrimmed returns a string representation of the error output of the process denoted
// by this struct with surrounding whitespace removed.
func (x RunResult) StderrTrimmed() string {
	return strings.TrimSpace(x.Stderr())
}

// Stderr returns a string representation of the stderr output of the process denoted
// by this struct.
func (x RunResult) Stderr() string {
	return x.stderrBuffer.String()
}

func (x RunResult) Error() error {
	return x.err
}

func (x RunResult) CombinedError() error {
	if x.err != nil {
		return x.err
	}

	if x.Successful() {
		return nil
	}
	return fmt.Errorf("execution of command '%s' failed: %s", x.command.String(), x.StderrTrimmed())
}

// StateString returns a string representation of the process denoted by
// this struct
func (x RunResult) StateString() string {
	state := x.ProcessState
	exitCode, err := x.ExitCode()
	exitCodeString := "?"
	if err == nil {
		exitCodeString = strconv.Itoa(exitCode)
	}
	return fmt.Sprintf("PID: %d, Exited: %t, Exit Code: %s, Success: %t, User Time: %s", state.Pid(), state.Exited(), exitCodeString, state.Success(), state.UserTime())
}

// ExitCode returns the exit code of the command denoted by this struct
func (x RunResult) ExitCode() (int, error) {
	var (
		waitStatus syscall.WaitStatus
		exitError  *exec.ExitError
	)
	ok := false
	if x.ProcessError != nil {
		exitError, ok = x.ProcessError.(*exec.ExitError)
	}
	if ok {
		waitStatus = exitError.Sys().(syscall.WaitStatus)
	} else {
		if x.ProcessState == nil {
			return -1, errors.New("no exit code available")
		}
		waitStatus = x.ProcessState.Sys().(syscall.WaitStatus)
	}
	return waitStatus.ExitStatus(), nil
}
