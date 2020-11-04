package tfcapture

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
	return &Capture{state}, nil
}

// CaptureOutputs writes all of the non-sensitive outputs to the io.Writer
// It returns the number of outputs written and err. If -1 is returned, there
// was an error in writing the outputs, and some outputs may have been written.
func (c *Capture) CaptureOutputs(w io.Writer) (int, error) {
	outputs := make(map[string]interface{})
	if c.state.Values == nil {
		return 0, nil
	}
	for name, output := range c.state.Values.Outputs {
		if !output.Sensitive {
			outputs[name] = output.Value
		}
	}
	if len(outputs) == 0 {
		return 0, nil
	}
	buf, err := json.MarshalIndent(outputs, "", "  ")
	if err != nil {
		return 0, err
	}
	_, err = w.Write(buf)
	if err != nil {
		return -1, err
	}
	return len(outputs), nil
}
