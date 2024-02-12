package component

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/jetbuild/runner/pkg/component"
	"github.com/tidwall/gjson"
)

type Context struct {
	arguments map[string]any
	outputs   *sync.Map
	index     uint
}

type TriggerContext struct {
	arguments    map[string]any
	OutputStream chan component.Output
}

func NewContext(arguments map[string]any, outputs *sync.Map, index uint) *Context {
	return &Context{
		arguments: arguments,
		outputs:   outputs,
		index:     index,
	}
}

func NewTriggerContext(arguments map[string]any) *TriggerContext {
	return &TriggerContext{
		arguments:    arguments,
		OutputStream: make(chan component.Output),
	}
}

func (c *Context) Output(o component.Output) {
	b, _ := json.Marshal(o)
	c.outputs.Store(c.index, gjson.ParseBytes(b))
}

func (t *TriggerContext) Output(o component.Output) {
	b, _ := json.Marshal(o)
	t.OutputStream <- gjson.ParseBytes(b)
}

func (c *Context) GetArgument(key string) any {
	argument, ok := c.arguments[key]
	if !ok {
		return nil
	}

	s := fmt.Sprint(argument)
	s, ok = strings.CutPrefix(s, "{{")
	if !ok {
		return argument
	}

	s, ok = strings.CutSuffix(s, "}}")
	if !ok {
		return argument
	}

	data := map[string]any{
		"outputs": make(map[string]any),
	}

	c.outputs.Range(func(key, value any) bool {
		data["outputs"].(map[string]any)[fmt.Sprint(key)] = value

		return true
	})

	m, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	return gjson.GetBytes(m, s).Value()
}

func (t *TriggerContext) GetArgument(key string) any {
	argument, ok := t.arguments[key]
	if !ok {
		return nil
	}

	return argument
}
