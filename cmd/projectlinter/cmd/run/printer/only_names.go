package printer

import (
	"fmt"
	"github.com/dozer111/projectlinter-core/util/painter"
	"strings"
)

// OnlyNames is mainly needed for debugging to see which rules was executed
type OnlyNames struct {
	painter painter.Painter
}

var _ Printer = (*OnlyNames)(nil)

func NewOnlyNames() *OnlyNames {
	return &OnlyNames{painter: painter.NewPainter()}
}

func (p *OnlyNames) Print(data []Data) string {
	result := make([]string, 0, 100)
	for _, set := range data {
		result = append(result, p.handleData(set)...)
	}

	return strings.Join(result, "\n")
}

func (p *OnlyNames) handleData(data Data) []string {
	setID := data.SetID
	if len(data.SetInitErrors) > 0 {
		return []string{p.painter.Red(fmt.Sprintf("== %s", setID))}
	}

	result := make([]string, 0, 50)
	result = append(result, p.painter.Yellow(fmt.Sprintf("== %s", setID)))
	for _, rule := range data.Rules {
		if rule.IsPassed() {
			result = append(
				result,
				fmt.Sprintf(
					"%s %s",
					p.painter.Yellow("="),
					p.painter.Green(fmt.Sprintf("✔ %s", rule.Title())),
				),
			)
			continue
		}

		result = append(
			result,
			fmt.Sprintf(
				"%s %s",
				p.painter.Yellow("="),
				p.painter.Red(fmt.Sprintf("✘ %s", rule.Title())),
			),
		)
	}

	result = append(result, "")
	return result
}
