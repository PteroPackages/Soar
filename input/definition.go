package input

type Node uint8

const (
	StringNode Node = iota
	NumberNode
	BoolNode
)

type Definition map[string]Node

func (d *Definition) Marshal(input map[string]string) (string, error)
