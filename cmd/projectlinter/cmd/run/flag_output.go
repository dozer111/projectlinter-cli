package run

import (
	"fmt"

	"github.com/spf13/pflag"
)

type Output string

const (
	OutputDefault Output = "default"
	OutputNames   Output = "names"
)

var _ pflag.Value = (*Output)(nil)

func (r *Output) Set(val string) error {
	newVal := Output(val)
	if !newVal.IsValid() {
		return fmt.Errorf("undefined output mode '%s'", val)
	}

	*r = newVal
	return nil
}

func (r *Output) Type() string {
	return fmt.Sprintf(`enum(%s|%s)`, OutputDefault, OutputNames)
}

func (r *Output) String() string {
	return string(*r)
}

func (r *Output) IsValid() bool {
	switch *r {
	case OutputDefault, OutputNames:
		return true
	}

	return false
}
