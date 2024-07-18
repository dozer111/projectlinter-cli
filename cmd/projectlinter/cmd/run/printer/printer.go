package printer

import "github.com/dozer111/projectlinter-core/rules"

type Printer interface {
	Print([]Data) string
}

type Data struct {
	SetID         string
	SetInitErrors []error
	Rules         []rules.Rule
}
