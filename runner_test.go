package gorun

import (
	"github.com/stretchr/testify/assert"
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

	cmd := LocalCommandFrom("python -c 'import os; os.write(1, b\"Yo\"); os.write(2, b\"def\")'")
	//.WithShell("sh", "-c")
	a.NotNil(cmd)
	res = NewWithCommand(cmd).
		WithoutStdout().
		WithoutStderr().
		LogCommand(false).
		MustExec()
	//fmt.Println("Stdout", res.StdoutTrimmed(), "Stderr", res.StderrTrimmed())
	a.Equal("Yo", res.Stdout())
	a.Equal("def", res.Stderr())

	// exit Fail
	cmd = LocalCommandFrom("python -c 'import sys; sys.exit(1)'")
	a.NotNil(cmd)
	res = NewWithCommand(cmd).
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
	err = res.CombinedError()
	a.Nil(err)
	a.True(res.Successful())

	// missing binary
	cmd = LocalCommandFrom("python-was-not-here")
	a.NotNil(cmd)
	res = NewWithCommand(cmd).
		LogCommand(false).
		MustExec()
	//fmt.Println("Stdout", res.StdoutTrimmed(), "Stderr", res.StderrTrimmed())
	a.True(res.Failed())
	a.NotNil(res.Error())
	err = res.CombinedError()
	a.NotNil(err)
	a.Equal("exec: \"python-was-not-here\": executable file not found in $PATH", err.Error())
}
