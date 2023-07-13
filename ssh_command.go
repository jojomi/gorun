package gorun

type SSHCommand struct {
	options       map[string]any
	remoteCommand Command
}
