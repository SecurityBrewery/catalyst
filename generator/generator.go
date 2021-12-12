package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/go-openapi/analysis"
	"github.com/go-swagger/go-swagger/generator"
	"github.com/iancoleman/strcase"
	"github.com/tidwall/sjson"
	"gopkg.in/yaml.v3"
)

//go:embed templates/simplemodel.gotmpl
var model embed.FS

func gotype(name string, s Schema, required []string) string {
	_, x := sgotype(name, s, required, false)
	return x
}

func sgotype(name string, s Schema, required []string, nopointer bool) (bool, string) {
	req := ""
	if !nopointer && !contains(required, name) {
		req = "*"
	}

	if s.Ref != "" {
		return false, req + path.Base(s.Ref)
	}

	primitive := false
	t := ""

	switch s.Type {
	case "string":
		if s.Format == "date-time" {
			t = req + "time.Time"
		} else {
			t = req + "string"
			primitive = true
		}
	case "boolean":
		t = req + "bool"
		primitive = true
	case "object":
		if s.AdditionalProperties != nil {
			subPrimitive, subType := sgotype(name, *s.AdditionalProperties, required, true)
			if subPrimitive {
				t = "map[string]" + subType
			} else {
				t = "map[string]*" + subType
			}
		} else {
			t = "interface{}"
		}
	case "number", "integer":
		if s.Format != "" {
			t = req + s.Format
		} else {
			t = req + "int"
		}
		primitive = true
	case "array":
		subPrimitive, subType := sgotype(name, *s.Items, required, true)
		if subPrimitive {
			t = "[]" + subType
		} else {
			t = "[]*" + subType
		}
	case "":
		t = "interface{}"
	default:
		panic(fmt.Sprintf("%#v", s))
	}

	return primitive, t
}

func omitempty(name string, required []string) bool {
	return !contains(required, name)
}

func contains(required []string, name string) bool {
	for _, r := range required {
		if r == name {
			return true
		}
	}
	return false
}

func tojson(name string, i Definition) string {
	b, _ := json.Marshal(i)
	b, _ = sjson.SetBytes(b, "$id", "#/definitions/"+name)
	return string(b)
}

func camel(s string) string {
	if s == "id" {
		return "ID"
	}
	return strcase.ToCamel(s)
}

func main() {
	flag.Parse()
	p := flag.Arg(0)

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	f, err := os.Open("generated/community.yml")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	s := Swagger{}
	dec := yaml.NewDecoder(f)
	err = dec.Decode(&s)
	if err != nil {
		log.Fatalln(err)
	}

	t := template.New("simplemodel.gotmpl")
	t.Funcs(map[string]interface{}{
		"camel":     camel,
		"gotype":    gotype,
		"omitempty": omitempty,
		"tojson":    tojson,
	})
	templ := template.Must(t.ParseFS(model, "templates/simplemodel.gotmpl"))

	err = os.MkdirAll("generated/models", os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	buf := bytes.NewBufferString("")

	props := map[string][]string{}
	for defName, definition := range s.Definitions {
		for propName := range definition.Properties {
			props[defName] = append(props[defName], propName)
		}
	}

	// for _, definition := range s.Definitions {
	// 	if definition.Embed != "" {
	// 		if parentProps, ok := props[definition.Embed]; ok {
	// 			for _, parentProp := range parentProps {
	// 				delete(definition.Properties, parentProp)
	// 			}
	// 		}
	// 	}
	// }

	err = templ.Execute(buf, &s)
	if err != nil {
		log.Fatalln(err)
	}

	fmtCode, err := format.Source(buf.Bytes())
	if err != nil {
		log.Println(err)
		fmtCode = buf.Bytes()
	}

	err = os.WriteFile("generated/models/models.go", fmtCode, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	generator.FuncMapFunc = func(opts *generator.LanguageOpts) template.FuncMap {
		df := generator.DefaultFuncMap(opts)

		df["path"] = func(basePath, lpath string, parameters generator.GenParameters) string {
			u := url.URL{Path: path.Join(basePath, lpath)}
			q := u.Query()
			for _, p := range parameters {
				if p.Location == "path" {
					if example, ok := p.Extensions["x-example"]; ok {
						u.Path = strings.ReplaceAll(u.Path, "{"+p.Name+"}", fmt.Sprint(example))
					}
				}
				if p.Location == "query" {
					if example, ok := p.Extensions["x-example"]; ok {
						q.Set(p.Name, fmt.Sprint(example))
					}
				}
			}
			u.RawQuery = q.Encode()
			return u.String()
		}
		df["body"] = func(parameters generator.GenParameters) interface{} {
			for _, p := range parameters {
				if p.Location == "body" {
					if example, ok := p.Extensions["x-example"]; ok {
						return example
					}
				}
			}
			return nil
		}
		df["ginizePath"] = func(path string) string {
			return strings.Replace(strings.Replace(path, "{", ":", -1), "}", "", -1)
		}
		df["export"] = func(name string) string {
			return strings.ToUpper(name[0:1]) + name[1:]
		}
		df["basePaths"] = func(operations []generator.GenOperation) []string {
			var l []string
			var seen = map[string]bool{}
			for _, operation := range operations {
				if _, ok := seen[operation.BasePath]; !ok {
					l = append(l, strings.TrimPrefix(operation.BasePath, "/"))
					seen[operation.BasePath] = true
				}
			}
			return l
		}
		df["roles"] = func(reqs []analysis.SecurityRequirement) string {
			for _, req := range reqs {
				if req.Name == "roles" {
					var roles []string
					for _, scope := range req.Scopes {
						roles = append(roles, "role."+strcase.ToCamel(strings.ReplaceAll(scope, ":", "_")))
						// roles = append(roles, permission.FromString(scope))
					}
					return strings.Join(roles, ", ")
				}
			}
			return ""
		}
		return df
	}

	opts := &generator.GenOpts{
		Spec:              "generated/community.yml",
		Target:            "generated",
		APIPackage:        "operations",
		ModelPackage:      "models",
		ServerPackage:     "restapi",
		ClientPackage:     "client",
		DefaultScheme:     "http",
		IncludeModel:      true,
		IncludeValidator:  true,
		IncludeHandler:    true,
		IncludeParameters: true,
		IncludeResponses:  true,
		IncludeURLBuilder: true,
		IncludeMain:       true,
		IncludeSupport:    true,
		ValidateSpec:      true,
		FlattenOpts: &analysis.FlattenOpts{
			Minimal: true,
			Verbose: true,
		},
		Name:              "catalyst-test",
		FlagStrategy:      "go-flags",
		CompatibilityMode: "modern",
		Sections: generator.SectionOpts{
			Application: []generator.TemplateOpts{
				{
					Name:     "api-server-test",
					Source:   path.Join(p, "templates/api_server_test.gotmpl"),
					Target:   "{{ .Target }}/test",
					FileName: "api_server_test.go",
				},
				// {
				// 	Name:       "configure",
				// 	Source:     "generator/config.gotmpl",
				// 	Target:     "{{ joinFilePath .Target .ServerPackage }}",
				// 	FileName:   "config.go",
				// 	SkipExists: false,
				// 	SkipFormat: false,
				// },
				{
					Name:     "embedded_spec",
					Source:   "asset:swaggerJsonEmbed",
					Target:   "{{ joinFilePath .Target .ServerPackage }}",
					FileName: "embedded_spec.go",
				},
				{
					Name:     "server",
					Source:   path.Join(p, "templates/api.gotmpl"),
					Target:   "{{ joinFilePath .Target .ServerPackage }}",
					FileName: "api.go",
				},
				{
					Name:     "response.go",
					Source:   path.Join(p, "templates/response.gotmpl"),
					Target:   "{{ .Target }}/restapi/api",
					FileName: "response.go",
				},
			},
			Operations: []generator.TemplateOpts{
				{
					Name:     "parameters",
					Source:   path.Join(p, "templates/parameter.gotmpl"),
					Target:   "{{ if gt (len .Tags) 0 }}{{ joinFilePath .Target .ServerPackage .APIPackage .Package  }}{{ else }}{{ joinFilePath .Target .ServerPackage .Package  }}{{ end }}",
					FileName: "{{ (snakize (pascalize .Name)) }}_parameters.go",
				},
			},
			Models: []generator.TemplateOpts{
				{
					Name:     "definition",
					Source:   "asset:model",
					Target:   "{{ joinFilePath .Target .ModelPackage }}/old",
					FileName: "{{ (snakize (pascalize .Name)) }}.go",
				},
				// {
				// 	Name:     "model",
				// 	Source:   "generator/model.gotmpl",
				// 	Target:   "{{ joinFilePath .Target .ModelPackage }}/old2",
				// 	FileName: "{{ (snakize (pascalize .Name)) }}.go",
				// },
			},
		},
	}

	err = opts.EnsureDefaults()
	if err != nil {
		log.Fatalln(err)
	}

	err = generator.GenerateServer("catalyst", nil, nil, opts)
	if err != nil {
		log.Fatalln(err)
	}
	// loads.Spec()
	// swagger.
}
