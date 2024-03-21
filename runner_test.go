package gorun

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestRunner(t *testing.T) {
	a := assert.New(t)

	res, err := New().
		WithCommand(LocalCommandFrom("echo abc")).
		WithoutStdout().
		Exec()
	a.Nil(err)
	a.True(res.Successful())
	a.Equal("abc", res.StdoutTrimmed())
}

func TestRunnerOutputStreams(t *testing.T) {
	a := assert.New(t)

	workingDir, err := os.Getwd()
	a.Nil(err)
	cmd := LocalCommandFrom("python3 -c 'import os, time; os.write(1, str.encode(\"lo\")); time.sleep(0.1); os.write(2, str.encode(os.getcwd()))'")
	a.NotNil(cmd)
	res := NewWithCommand(cmd).
		WithoutStdout().
		WithoutStderr().
		InWorkingDir(workingDir).
		MustExec()
	a.Equal("lo", res.Stdout())
	a.Equal(workingDir, res.Stderr())
	a.Equal("lo"+workingDir, res.CombinedOutput())
}

func TestRunnerEnv(t *testing.T) {
	a := assert.New(t)

	cmd := LocalCommandFrom("python3 -c 'import os; os.write(1, str.encode(os.environ[\\'GORUNTEST\\'] + \"lo\")); os.write(2, str.encode(os.environ[\\'PATH\\']));'")
	a.NotNil(cmd)
	res := NewWithCommand(cmd).
		AddEnv("GORUNTEST", "Yo").
		WithoutStdout().
		WithoutStderr().
		MustExec()
	a.True(res.Successful())
	a.Equal("Yolo", res.Stdout())
	a.NotEqual("", res.StderrTrimmed())
}

func TestRunnerReset(t *testing.T) {
	a := assert.New(t)

	cmd := LocalCommandFrom("echo 'abc'")
	a.NotNil(cmd)
	runner := NewWithCommand(cmd).
		WithoutStdout().
		WithoutStderr()
	res := runner.MustExec()

	res = runner.Reset().WithCommand(LocalCommandFrom("echo 'dev'")).MustExec()

	a.True(res.Successful())
	a.Equal("dev", res.StdoutTrimmed())
	a.Equal("", res.StderrTrimmed())
}

func TestExitCodes(t *testing.T) {
	a := assert.New(t)

	// exit Fail
	cmd := LocalCommandFrom("python3 -c 'import sys; sys.exit(1)'")
	a.NotNil(cmd)
	res := NewWithCommand(cmd).
		MustExec()
	a.True(res.Failed())
	a.False(res.Successful())

	// NonZeroExitOK
	res = NewWithCommand(cmd).
		NonZeroExitOK().
		MustExec()
	a.False(res.Failed())
	err := res.CombinedError()
	a.Nil(err)
	a.True(res.Successful())
}

func TestInvalidCommand(t *testing.T) {
	a := assert.New(t)

	// missing binary
	cmd := LocalCommandFrom("python-was-not-here")
	a.NotNil(cmd)
	res := NewWithCommand(cmd).
		MustExec()
	a.True(res.Failed())
	a.NotNil(res.Error())
	err := res.CombinedError()
	a.NotNil(err)
	a.Equal("exec: \"python-was-not-here\": executable file not found in $PATH", err.Error())
}
