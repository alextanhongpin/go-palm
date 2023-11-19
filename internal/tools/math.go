package tools

import (
	"fmt"
	"strings"

	"github.com/antonmedv/expr"
)

type Math struct {
	name        string
	tag         string
	description string
}

func NewMath() *Math {
	return &Math{
		name: "Calculator",
		tag:  "calc",
		description: `When solving this problem, use the calculator for any arithmetic.

To use the calculator, put an expression between <calc></calc> tags.
The answer will be printed after the </calc> tag.

For example: 2 houses  * 8 cats/house = <calc>2 * 8</calc> = 16 cats`,
	}
}

func (m *Math) Name() string {
	return m.name
}

func (m *Math) Description() string {
	return m.description
}

func (m *Math) Tag() string {
	return m.tag
}

func (m *Math) Eval(description string) (string, error) {
	startTag := fmt.Sprintf("<%s>", m.tag)
	endTag := fmt.Sprintf("</%s>", m.tag)
	p, e, ok := strings.Cut(description, startTag)
	if !ok {
		return p, nil
	}

	v, err := m.calc(e)
	if err != nil {
		return p, err
	}

	endPrompt := fmt.Sprintf("%s%s%s%s = %v", p, startTag, e, endTag, v)

	return endPrompt, nil
}

func (m *Math) calc(code string) (any, error) {
	program, err := expr.Compile(code, expr.Env(nil))
	if err != nil {
		return nil, err
	}

	output, err := expr.Run(program, nil)
	if err != nil {
		return nil, err
	}

	return output, nil
}
