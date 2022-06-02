//go:build ignore

package main

import (
	"os"
	"regexp"
	"strings"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/iancoleman/strcase"
	"github.com/leg100/otf"
	"github.com/leg100/otf/http/html"
)

var (
	varRx = regexp.MustCompile(`\{([^\:\}]+)\}`)
)

// Dummy app to pass in
type app struct {
	otf.Application
}

// helper is the path helper to be generated
type helper struct {
	// name and path contains every needed to generate a path helper
	name, path string
	// func params extracted from path
	params []string
}

func newHelper(name, path string) helper {
	return helper{
		name:   name,
		path:   path,
		params: parsePathVars(path),
	}
}

// funcName returns the name of the helper function
func (h helper) funcName() string {
	return h.name + "Path"
}

// return the function definition as a string
func (h helper) funcDefStr() string {
	b := strings.Builder{}
	b.WriteString("func ")
	b.WriteString(h.funcName())
	b.WriteRune('(')
	b.WriteString(strings.Join(h.params, ", "))
	if len(h.params) == 0 {
		b.WriteString(") string {\n")
		b.WriteString("\treturn \"")
		b.WriteString(h.fmtStr())
		b.WriteString("\"\n")
	} else {
		// parameters share a string type
		b.WriteString(" string")
		b.WriteString(") string {\n")
		b.WriteString("\treturn fmt.Sprintf(\"")
		b.WriteString(h.fmtStr())
		b.WriteString("\", ")
		// fmt.Sprintf params
		b.WriteString(strings.Join(h.params, ", "))
		b.WriteString(")\n")
	}
	b.WriteString("}\n\n")
	return b.String()
}

// funcMapAssignmentStr returns as a string the assignment of the helper
// function to a template func map, i.e.:
//
//     m["getWorkspacePath"] = getWorkspacePath
func (h helper) funcMapAssignmentStr() string {
	b := strings.Builder{}
	b.WriteString("\tm[\"")
	b.WriteString(h.funcName())
	b.WriteString("\"] = ")
	b.WriteString(h.funcName())
	b.WriteRune('\n')
	return b.String()
}

func (h helper) fmtStr() string {
	return varRx.ReplaceAllString(h.path, "%s")
}

// parsePathVars parses a mux route's path variables. It also ensures they're in
// lower camel case, suitable for unexported go variables.
func parsePathVars(path string) []string {
	matches := varRx.FindAllStringSubmatch(path, -1)
	if matches == nil {
		return nil
	}
	var vars []string
	for _, m := range matches {
		vars = append(vars, strcase.ToLowerCamel(m[1]))
	}
	return vars
}

func main() {
	// get routes from web app
	r := mux.NewRouter()
	html.AddRoutes(logr.Discard(), html.Config{}, &app{}, r)
	// walk routes and populate helpers
	var helpers []helper
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		if route.GetName() == "" {
			// skip routes without a name
			return nil
		}
		path, err := route.GetPathTemplate()
		if err != nil {
			// skip routes without a path
			return nil
		}
		helpers = append(helpers, newHelper(route.GetName(), path))
		return nil
	})
	// build str for output and write to file
	b := strings.Builder{}
	b.WriteString("package html\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"fmt\"\n")
	b.WriteString("\t\"html/template\"\n")
	b.WriteString(")\n\n")
	for _, h := range helpers {
		b.WriteString(h.funcDefStr())
	}
	b.WriteString("func addHelpersToFuncMap(m template.FuncMap) {\n")
	for _, h := range helpers {
		b.WriteString(h.funcMapAssignmentStr())
	}
	b.WriteString("}\n")
	os.WriteFile("path_helpers_gen.go", []byte(b.String()), 0644)
}
