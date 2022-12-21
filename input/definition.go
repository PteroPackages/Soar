package input

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Node uint8

const (
	StringNode Node = iota
	ArrayStringNode
	NullStringNode
	NumberNode
	BoolNode
)

type Definition map[string]Node

func Marshal(def Definition, input map[string]string) ([]byte, error) {
	p := map[string]interface{}{}

	for k, n := range def {
		v, ok := input[k]
		if !ok {
			return nil, errors.New("missing key in input for definition")
		}

		switch n {
		case StringNode:
			p[k] = v
		case ArrayStringNode:
			p[k] = strings.Split(v, ",")
		case NullStringNode:
			if v == "null" {
				p[k] = nil
			} else {
				p[k] = v
			}
		case NumberNode:
			r, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid integer \"%s\"", v)
			}

			p[k] = r
		case BoolNode:
			r, err := strconv.ParseBool(v)
			if err != nil {
				return nil, fmt.Errorf("invalid boolean: \"%s\"", v)
			}

			p[k] = r
		}
	}

	return json.Marshal(p)
}
