package tools

import (
	"bytes"
	"text/template"
)

type Template string

func (t Template) Format(data any) string {
	tmpl := template.Must(template.New("").Parse(string(t)))

	var b bytes.Buffer
	err := tmpl.Execute(&b, data)
	if err != nil {
		panic(err)
	}

	return string(b.Bytes())
}
