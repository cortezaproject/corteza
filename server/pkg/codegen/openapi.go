package codegen

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// OpenAPI documents for the Swagger UI the server serves at /docs, generated
// from the same rest.yaml files the handlers are generated from, so the page
// cannot describe a route the server does not have.
//
// rest.yaml does not describe response bodies, so every operation declares one
// generic 200: the API answers {"response": …} on success and {"error": …}
// on failure, both with HTTP 200, and the docs page says so once.

type (
	oaDoc struct {
		OpenAPI string                      `yaml:"openapi"`
		Info    oaInfo                      `yaml:"info"`
		Servers []oaServer                  `yaml:"servers"`
		Tags    []oaTag                     `yaml:"tags"`
		Paths   map[string]map[string]*oaOp `yaml:"paths"`
	}
	oaInfo struct {
		Title       string `yaml:"title"`
		Description string `yaml:"description"`
		Version     string `yaml:"version"`
	}
	oaServer struct {
		URL string `yaml:"url"`
	}
	oaTag struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description,omitempty"`
	}
	oaOp struct {
		OperationID string            `yaml:"operationId"`
		Summary     string            `yaml:"summary,omitempty"`
		Tags        []string          `yaml:"tags"`
		Parameters  []oaParam         `yaml:"parameters,omitempty"`
		RequestBody *oaBody           `yaml:"requestBody,omitempty"`
		Responses   map[string]oaResp `yaml:"responses"`
	}
	oaParam struct {
		Name        string   `yaml:"name"`
		In          string   `yaml:"in"`
		Required    bool     `yaml:"required"`
		Description string   `yaml:"description,omitempty"`
		Schema      oaSchema `yaml:"schema"`
	}
	oaBody struct {
		Required bool               `yaml:"required"`
		Content  map[string]oaMedia `yaml:"content"`
	}
	oaMedia struct {
		Schema oaObject `yaml:"schema"`
	}
	oaObject struct {
		Type       string              `yaml:"type"`
		Properties map[string]oaSchema `yaml:"properties,omitempty"`
		Required   []string            `yaml:"required,omitempty"`
	}
	oaSchema struct {
		Type        string    `yaml:"type,omitempty"`
		Format      string    `yaml:"format,omitempty"`
		Items       *oaSchema `yaml:"items,omitempty"`
		Description string    `yaml:"description,omitempty"`
	}
	oaResp struct {
		Description string             `yaml:"description"`
		Content     map[string]oaMedia `yaml:"content,omitempty"`
	}
)

const oaDescription = "Generated from server/%s/rest.yaml by make codegen; edit that file, not this one. " +
	"Every call answers HTTP 200: a success body is {\"response\": …} and a failure is {\"error\": {\"message\": …}}, " +
	"so read the body, not the status. IDs are 64-bit integers carried as strings. " +
	"Collection paths end in a slash; a compose update is a POST to the resource path."

// openAPIDoc builds the document for one app.
func openAPIDoc(d *restDef) *oaDoc {
	app := path.Base(d.outputDir)
	doc := &oaDoc{
		OpenAPI: "3.0.3",
		Info: oaInfo{
			Title:       "Human " + export(app) + " API",
			Description: fmt.Sprintf(oaDescription, app),
			Version:     "generated",
		},
		Servers: []oaServer{{URL: "/api"}},
		Paths:   map[string]map[string]*oaOp{},
	}

	for _, e := range d.Endpoints {
		doc.Tags = append(doc.Tags, oaTag{Name: e.Title, Description: strings.TrimSpace(e.Description)})
		for _, a := range e.Apis {
			p := "/" + app + e.Path + a.Path
			if doc.Paths[p] == nil {
				doc.Paths[p] = map[string]*oaOp{}
			}
			doc.Paths[p][strings.ToLower(a.Method)] = openAPIOp(e, a)
		}
	}
	return doc
}

func openAPIOp(e *restEndpointDef, a *restEndpointApi) *oaOp {
	op := &oaOp{
		OperationID: e.Entrypoint + "." + a.Name,
		Summary:     a.Title,
		Tags:        []string{e.Title},
		Responses: map[string]oaResp{
			"200": {
				Description: "The result under \"response\", or the failure under \"error\".",
				Content:     map[string]oaMedia{"application/json": {Schema: oaObject{Type: "object"}}},
			},
		},
	}

	// Endpoint-level path params were merged into each api by procRest, so
	// a.Params carries everything, in the order the yaml declares.
	for _, p := range a.Params.Path {
		op.Parameters = append(op.Parameters, oaParam{Name: p.Name, In: "path", Required: true, Description: p.Title, Schema: openAPISchema(p.Type)})
	}
	for _, p := range a.Params.Get {
		op.Parameters = append(op.Parameters, oaParam{Name: p.Name, In: "query", Required: p.Required, Description: p.Title, Schema: openAPISchema(p.Type)})
	}

	if len(a.Params.Post) > 0 {
		body := oaObject{Type: "object", Properties: map[string]oaSchema{}}
		media := "application/json"
		for _, p := range a.Params.Post {
			s := openAPISchema(p.Type)
			s.Description = p.Title
			body.Properties[p.Name] = s
			if p.Required {
				body.Required = append(body.Required, p.Name)
			}
			if p.IsUpload() {
				media = "multipart/form-data"
			}
		}
		sort.Strings(body.Required)
		op.RequestBody = &oaBody{Required: len(body.Required) > 0, Content: map[string]oaMedia{media: {Schema: body}}}
	}
	return op
}

// openAPISchema maps a rest.yaml Go type onto a wire type. IDs travel as
// strings; anything domain-shaped is an object the docs describe elsewhere.
func openAPISchema(goType string) oaSchema {
	t := strings.TrimPrefix(goType, "*")
	if strings.HasPrefix(t, "[]") {
		item := openAPISchema(t[2:])
		return oaSchema{Type: "array", Items: &item}
	}
	switch t {
	case "string":
		return oaSchema{Type: "string"}
	case "uint64", "int64":
		return oaSchema{Type: "string", Format: t, Description: "64-bit integer as a string"}
	case "uint", "int", "uint32", "int32":
		return oaSchema{Type: "integer"}
	case "float", "float64", "float32":
		return oaSchema{Type: "number"}
	case "bool":
		return oaSchema{Type: "boolean"}
	case "time.Time":
		return oaSchema{Type: "string", Format: "date-time"}
	case "multipart.FileHeader":
		return oaSchema{Type: "string", Format: "binary"}
	}
	if strings.HasPrefix(t, "map[") || strings.Contains(t, ".") || strings.HasSuffix(t, "Set") {
		return oaSchema{Type: "object", Description: t}
	}
	return oaSchema{Type: "string", Description: t}
}

// renderOpenAPI is the YAML text for one app's document.
func renderOpenAPI(d *restDef) ([]byte, error) {
	var b strings.Builder
	b.WriteString("# Generated by make codegen from " + path.Base(d.outputDir) + "/rest.yaml. Do not edit.\n")
	out, err := yaml.Marshal(openAPIDoc(d))
	if err != nil {
		return nil, err
	}
	b.Write(out)
	return []byte(b.String()), nil
}

// openAPIPath is where the server's /docs page reads the document from.
func openAPIPath(d *restDef) string {
	return path.Join(path.Dir(d.outputDir), "docs", path.Base(d.outputDir)+".yaml")
}

func genOpenAPI(dd ...*restDef) error {
	for _, d := range dd {
		out, err := renderOpenAPI(d)
		if err != nil {
			return fmt.Errorf("openapi for %s: %w", d.App, err)
		}
		if err := os.WriteFile(openAPIPath(d), out, 0o644); err != nil {
			return err
		}
	}
	return nil
}
