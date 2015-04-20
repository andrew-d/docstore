package docstore

import (
	"encoding/json"
	"text/template"
)

var (
	funcMap = template.FuncMap{
		"json": templateToJson,
	}
)

func ParseFormat(s string) (*template.Template, error) {
	return template.New("format").Funcs(funcMap).Parse(s)
}

func templateToJson(input interface{}) (string, error) {
	by, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	return string(by), nil
}
