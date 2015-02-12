package terraform

import (
	"fmt"
)

// EvalWriteOutput is an EvalNode implementation that writes the output
// for the given name to the current state.
type EvalWriteOutput struct {
	Name  string
	Value EvalNode
}

func (n *EvalWriteOutput) Args() ([]EvalNode, []EvalType) {
	return []EvalNode{n.Value}, []EvalType{EvalTypeConfig}
}

// TODO: test
func (n *EvalWriteOutput) Eval(
	ctx EvalContext, args []interface{}) (interface{}, error) {
	config := args[0].(*ResourceConfig)

	state, lock := ctx.State()
	if state == nil {
		return nil, fmt.Errorf("cannot write state to nil state")
	}

	// Get a write lock so we can access this instance
	lock.Lock()
	defer lock.Unlock()

	// Look for the module state. If we don't have one, create it.
	mod := state.ModuleByPath(ctx.Path())
	if mod == nil {
		mod = state.AddModule(ctx.Path())
	}

	// Get the value from the config
	valueRaw, ok := config.Get("value")
	if !ok {
		valueRaw = ""
	}

	// If it is a list of values, get the first one
	if list, ok := valueRaw.([]interface{}); ok {
		valueRaw = list[0]
	}
	if _, ok := valueRaw.(string); !ok {
		valueRaw = ""
	}

	// Write the output
	mod.Outputs[n.Name] = valueRaw.(string)

	return nil, nil
}

func (n *EvalWriteOutput) Type() EvalType {
	return EvalTypeNull
}
