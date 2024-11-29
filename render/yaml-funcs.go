// There are some functions like in helm
package render

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v3"
)

// tpl function to render a template with data, supporting recursion
func tpl(tplStr string, data interface{}, funcMap template.FuncMap) (string, error) {
	tmpl, err := template.New("tpl").Funcs(funcMap).Parse(tplStr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func funcMap() template.FuncMap {
	f := sprig.TxtFuncMap()
	delete(f, "env")
	delete(f, "expandenv")
	var extra template.FuncMap
	extra = template.FuncMap{
		"toYaml":   toYAML,
		"toJson":   toJSON,
		"fromYaml": fromYAML,
		"tpl": func(tplStr string, data interface{}) (string, error) {
			return tpl(tplStr, data, extra)
		},

		// functions are not implemented and I don't want to
		"include":  func(string, interface{}) string { return "not implemented" },
		"required": func(string, interface{}) string { return "not implemented" },
	}

	for k, v := range extra {
		f[k] = v
	}

	return f
}

func toYAML(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		log.Fatalf("Error: Can't execute toYAML func:\"%v\"\n   \"%s\"", err, v)
	}
	return strings.TrimSuffix(string(data), "\n")
}

func toJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("Error: Can't execute toJSON func:\"%v\"\n   \"%s\"", err, v)
	}
	return string(data)
}

func fromYAML(str string) map[string]interface{} {
	m := map[string]interface{}{}

	if err := yaml.Unmarshal([]byte(str), &m); err != nil {
		m["Error"] = err.Error()
	}
	return m
}
