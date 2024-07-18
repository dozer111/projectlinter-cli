package printer

import (
	"fmt"
	"github.com/dozer111/projectlinter-core/rules"
	"github.com/dozer111/projectlinter-core/util/painter"
	"strings"
)

// Default - default verbose output which show maximum information
type Default struct {
	painter painter.Painter
}

var _ Printer = (*Default)(nil)

func NewDefault() *Default {
	return &Default{painter: painter.NewPainter()}
}

func (p *Default) Print(data []Data) string {
	header := make([]string, 0, len(data))
	errorsPart := make([]string, 0, 128)
	for _, setData := range data {
		setBrief, setErrors := p.handleData(setData)
		header = append(header, setBrief)
		errorsPart = append(errorsPart, setErrors...)
	}

	result := []string{
		strings.Join(header, ""),
		"",
	}
	result = append(result, errorsPart...)

	if len(errorsPart) == 0 {
		result = append(
			result,
			"Well done! All rules are passed!",
		)
	}

	return strings.Join(result, "\n")
}

func (p *Default) handleData(data Data) (string, []string) {
	setID := data.SetID
	if len(data.SetInitErrors) > 0 {
		return "", p.addSetErrors(data)
	}

	var failedRules []rules.Rule
	brief := ""
	for _, rule := range data.Rules {
		if rule.IsPassed() {
			brief += "."
			continue
		}

		brief += "F"
		failedRules = append(failedRules, rule)
	}

	return brief, p.addFailedRules(setID, failedRules)
}

func (p *Default) addSetErrors(data Data) []string {
	result := make([]string, 0, 16)

	result = append(
		result,
		fmt.Sprintf("== %s\n", data.SetID),
		" Cannot run set because of errors:\n",
	)

	for _, e := range data.SetInitErrors {
		result = append(result, fmt.Sprintf("\t- %s", e))
	}

	return result
}

func (p *Default) addFailedRules(setID string, failedRules []rules.Rule) []string {
	if len(failedRules) == 0 {
		return []string{}
	}

	result := make([]string, 0, 64)

	result = append(result, fmt.Sprintf("== %s", setID))

	for _, r := range failedRules {
		result = append(result, p.painter.Yellow(fmt.Sprintf("┌%s"+
			"\n│ > %s\n│ id: %s\n│", strings.Repeat("-", 100), r.Title(), r.ID())))

		failedMessage := r.FailedMessage()
		failedMessageHasMultilineSolution := false

		// if solution has more than 1 line of code to add
		// we need to skip the | in the start of line to easy copy-paste
		{
			solutionCounter := 0
			for _, o := range failedMessage {
				if solutionCounter > 1 {
					failedMessageHasMultilineSolution = true
					break
				}
				// string contains green text
				if strings.Contains(o, "\u001B[32m") {
					solutionCounter++
				}
			}
		}

		for _, o := range failedMessage {
			message := fmt.Sprintf("%s   %s", p.painter.Yellow("│"), o)
			// string contains green text
			if strings.Contains(o, "\u001B[32m") && failedMessageHasMultilineSolution {
				message = fmt.Sprintf("    %s", o)
			}

			result = append(result, message)
		}

		result = append(result, p.painter.Yellow("└"+strings.Repeat("-", 100)))
	}
	result = append(result, "")

	return result
}
