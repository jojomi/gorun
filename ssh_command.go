package gorun

import "fmt"

type SSHCommand struct {
	alias         string
	options       map[string]any
	remoteCommand Command
}

func NewSSHCommandFrom(alias string, command Command) *SSHCommand {
	cmd := &SSHCommand{
		alias:         alias,
		remoteCommand: command,
	}
	return cmd
}

func (x SSHCommand) String() string {
	return fmt.Sprintf(`ssh "%s" -- %s`, x.alias, x.remoteCommand.String())
}

func (x SSHCommand) Binary() string {
	return "ssh"
}

func (x SSHCommand) Args() []string {
	args := []string{
		x.alias,
		`--`,
		x.remoteCommand.Binary(),
	}

	args = append(args, x.remoteCommand.Args()...)
	return args
}
