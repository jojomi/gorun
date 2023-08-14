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
	//fmt.Println("Stdout", res.StdoutTrimmed(), "Stderr", res.StderrTrimmed())
	a.True(res.Successful())
	a.Equal("abc", res.StdoutTrimmed())
}

func TestRunnerOutputStreamsAndEnv(t *testing.T) {
	a := assert.New(t)

	workingDir, err := os.Getwd()
	a.Nil(err)
	cmd := LocalCommandFrom("python3 -c 'import os; os.write(1, str.encode(os.environ[\\'GORUNTEST\\'] + \"lo\")); os.write(2, str.encode(os.getcwd()))'")
	//.WithShell("sh", "-c")
	a.NotNil(cmd)
	res := NewWithCommand(cmd).
		AddEnv("GORUNTEST", "Yo").
		WithoutStdout().
		WithoutStderr().
		InWorkingDir(workingDir).
		LogCommand(true).
		MustExec()
	//fmt.Println("Stdout", res.StdoutTrimmed(), "Stderr", res.StderrTrimmed())
	a.Equal("Yolo", res.Stdout())
	a.Equal(workingDir, res.Stderr())
	a.Equal("Yolo"+workingDir, res.CombinedOutput())
}

func TestExitCodes(t *testing.T) {
	a := assert.New(t)

	// exit Fail
	cmd := LocalCommandFrom("python3 -c 'import sys; sys.exit(1)'")
	a.NotNil(cmd)
	res := NewWithCommand(cmd).
		LogCommand(false).
		MustExec()
	//fmt.Println("Stdout", res.StdoutTrimmed(), "Stderr", res.StderrTrimmed())
	a.True(res.Failed())
	a.False(res.Successful())

	// NonZeroExitOK
	res = NewWithCommand(cmd).
		NonZeroExitOK().
		LogCommand(false).
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
		LogCommand(false).
		MustExec()
	//fmt.Println("Stdout", res.StdoutTrimmed(), "Stderr", res.StderrTrimmed())
	a.True(res.Failed())
	a.NotNil(res.Error())
	err := res.CombinedError()
	a.NotNil(err)
	a.Equal("exec: \"python-was-not-here\": executable file not found in $PATH", err.Error())
}
