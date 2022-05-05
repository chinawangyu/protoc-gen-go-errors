package main

import (
	"bytes"
	"text/template"
)

var errorsTemplate = `
{{ range .Errors }}

func Is{{.CamelValue}}(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	stat := e.Metadata["stat"]
	return (e.Reason == {{.Name}}_{{.Value}}.String() || stat == "{{.StatCode}}") && e.Code == {{.HTTPCode}}
}

func Error{{.CamelValue}}(format string, args ...interface{}) *errors.Error {
	md := map[string]string{
		"stat": "{{.StatCode}}",
	}
	 return errors.New({{.HTTPCode}}, {{.Name}}_{{.Value}}.String(), fmt.Sprintf(format, args...)).WithMetadata(md)
}

func ErrorMsg{{.CamelValue}}(args ...interface{}) *errors.Error {
	md := map[string]string{
		"stat": "{{.StatCode}}",
	}

	{{if ne .ErrMsg ""}}
	return errors.New({{.HTTPCode}}, {{.Name}}_{{.Value}}.String(), fmt.Sprintf("{{.ErrMsg}}", args...)).WithMetadata(md)
	{{else}}
	return errors.New({{.HTTPCode}}, {{.Name}}_{{.Value}}.String(), fmt.Sprintf("服务异常")).WithMetadata(md)
	{{end}}
}

{{- end }}
`

type errorInfo struct {
	Name       string
	Value      string
	HTTPCode   int
	StatCode   int
	CamelValue string
	ErrMsg     string
}

type errorWrapper struct {
	Errors []*errorInfo
}

func (e *errorWrapper) execute() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("errors").Parse(errorsTemplate)
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, e); err != nil {
		panic(err)
	}
	return buf.String()
}
