package tfplananizer

import (
	"encoding/json"
	"io"

	tfjson "github.com/hashicorp/terraform-json"
)

type Capture struct {
	state *tfjson.State
}

func NewFromState(r io.Reader) (*Capture, error) {
	var state *tfjson.State
	if err := json.NewDecoder(r).Decode(&state); err != nil {
		return nil, err
	}

	if err := state.Validate(); err != nil {
		return nil, err
	}
	return &Capture{state}, nil
}

// CaptureOutputs writes all of the non-sensitive outputs to the io.Writer
func (c *Capture) CaptureOutputs(w io.Writer) error {
	outputs := make(map[string]interface{})
	for name, output := range c.state.Values.Outputs {
		if !output.Sensitive {
			outputs[name] = output.Value
		}
	}
	if len(outputs) == 0 {
		return nil
	}
	buf, err := json.Marshal(outputs)
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}
