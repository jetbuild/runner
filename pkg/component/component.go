package component

type Context interface {
	Output(Output)
	GetArgument(key string) any
}

type Output any
