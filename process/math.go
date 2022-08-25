package process

import (
	"context"
	"fmt"

	"github.com/brexhq/substation/condition"
	"github.com/brexhq/substation/config"
)

/*
Math processes data by applying mathematic operations. The processor supports these patterns:
	JSON:
		{"math":[1,3]} >>> {"math":4}

When loaded with a factory, the processor uses this JSON configuration:
	{
		"type": "math",
		"settings": {
			"options": {
				"operation": "add"
			},
			"input_key": "math",
			"output_key": "math"
		}
	}
*/
type Math struct {
	Options   MathOptions      `json:"options"`
	Condition condition.Config `json:"condition"`
	InputKey  string           `json:"input_key"`
	OutputKey string           `json:"output_key"`
}

/*
MathOptions contains custom options for the Math processor:
	Operation:
		the operator applied to the data
		must be one of:
			add
			subtract
			multiply
			divide
*/
type MathOptions struct {
	Operation string `json:"operation"`
}

// ApplyBatch processes a slice of encapsulated data with the Math processor. Conditions are optionally applied to the data to enable processing.
func (p Math) ApplyBatch(ctx context.Context, caps []config.Capsule) ([]config.Capsule, error) {
	op, err := condition.OperatorFactory(p.Condition)
	if err != nil {
		return nil, fmt.Errorf("process math applybatch: %v", err)
	}

	caps, err = conditionallyApplyBatch(ctx, caps, op, p)
	if err != nil {
		return nil, fmt.Errorf("process math applybatch: %v", err)
	}

	return caps, nil
}

// Apply processes encapsulated data with the Math processor.
func (p Math) Apply(ctx context.Context, cap config.Capsule) (config.Capsule, error) {
	// error early if required options are missing
	if p.Options.Operation == "" {
		return cap, fmt.Errorf("process math apply: options %+v: %v", p.Options, ProcessorMissingRequiredOptions)
	}

	// only supports JSON, error early if there are no keys
	if p.InputKey == "" && p.OutputKey == "" {
		return cap, fmt.Errorf("process math apply: inputkey %s outputkey %s: %v", p.InputKey, p.OutputKey, ProcessorInvalidDataPattern)
	}

	var value int64
	result := cap.Get(p.InputKey)
	for i, res := range result.Array() {
		if i == 0 {
			value = res.Int()
			continue
		}

		switch p.Options.Operation {
		case "add":
			value += res.Int()
		case "subtract":
			value -= res.Int()
		case "multiply":
			value = value * res.Int()
		case "divide":
			value = value / res.Int()
		}
	}

	if err := cap.Set(p.OutputKey, value); err != nil {
		return cap, fmt.Errorf("process math apply: %v", err)
	}

	return cap, nil
}
