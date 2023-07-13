package gorun

type Command interface {
	String() string
	Binary() string
	Args() []string
}
