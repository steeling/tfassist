package tfplananizer

import (
	"encoding/json"
	"io"

	tfjson "github.com/hashicorp/terraform-json"
)

type Plananizer struct {
	plan *tfjson.Plan
}

func NewFromPlan(r io.Reader) (*Plananizer, error) {
	var plan *tfjson.Plan
	if err := json.NewDecoder(r).Decode(&plan); err != nil {
		return nil, err
	}
	return &Plananizer{plan}, nil
}

// Returns true/false if anything is destroyed, and the addresses of destroyed resources.
func (p *Plananizer) Destroys() (bool, []string) {
	var resources []string
	for _, rc := range p.plan.ResourceChanges {
		if rc.Change != nil && rc.Change.Actions.Delete() {
			resources = append(resources, rc.Address)
		}
	}
	return len(resources) > 0, resources
}
