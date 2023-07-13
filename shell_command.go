package gorun

type ShellCommand struct {
	LocalCommand

	shellPath    string
	shellArgs    string
	innerCommand Command
}

func NewShellCommandFrom(command Command) *ShellCommand {
	cmd := &ShellCommand{
		innerCommand: command,
	}
	return cmd.WithShell("sh", "-c")
}

func LocalShellCommandFrom(command string) *ShellCommand {
	return NewShellCommandFrom(LocalCommandFrom(command))
}

func (x *ShellCommand) WithShell(shellPath string, shellArgs string) *ShellCommand {
	x.shellPath = shellPath
	x.shellArgs = shellArgs
	return x
}

func (x *ShellCommand) String() string {
	return NewLocalCommand().AddAll(x.shellPath, x.shellArgs, x.innerCommand.String()).String()
}
