package html

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"regexp"
	"sort"

	"github.com/egorbanin/speka/speka"
	"github.com/hjson/hjson-go/v4"
)

//go:embed tpl/*.html
var tplFS embed.FS

type Generator struct {
	tpl *template.Template
}

func NewGenerator(tpl *template.Template) *Generator {
	return &Generator{tpl: tpl}
}

func CreateGenerator() (*Generator, error) {
	t, err := template.ParseFS(tplFS, "tpl/*.html")
	if err != nil {
		return nil, err
	}

	return NewGenerator(t), nil
}

func (g *Generator) Html(s speka.Speka, w io.Writer) error {
	names := make([]string, 0, len(s.Methods))
	for name := range s.Methods {
		names = append(names, name)
	}
	sort.Strings(names)

	methods := make([]methodView, 0, len(names))
	menu := make(menu, 0, len(names))
	for _, name := range names {
		m := s.Methods[name]
		id := methodID(name)

		params, err := nodeFields(name, "rq", m.Rq)
		if err != nil {
			return err
		}
		response, err := nodeFields(name, "rs", m.Rs)
		if err != nil {
			return err
		}

		methods = append(methods, methodView{
			ID:       id,
			Name:     name,
			Desc:     m.Desc,
			Params:   params,
			Response: response,
			Errors:   m.Errors,
		})
		menu = append(menu, menuItem{
			Text: name,
			Link: id,
		})
	}

	return g.tpl.ExecuteTemplate(w, "layout.html", layout{
		Name:    s.Name,
		Desc:    s.Desc,
		Menu:    menu,
		Methods: methods,
	})
}

func nodeFields(name, part string, node *hjson.Node) ([]*speka.Property, error) {
	if node == nil {
		return nil, nil
	}

	p, err := speka.ParseProperty(part, node)
	if err != nil {
		return nil, fmt.Errorf("speka.ParseProperty %s %s: %w", name, part, err)
	}

	if p.Kind == speka.KindObject {
		return p.Properties, nil
	}

	return []*speka.Property{p}, nil
}

var (
	nonID      = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)
	trimDashes = regexp.MustCompile(`^-+|-+$`)
)

func methodID(name string) string {
	id := trimDashes.ReplaceAllString(nonID.ReplaceAllString(name, "-"), "")
	if id == "" {
		return "method"
	}
	if id[0] >= '0' && id[0] <= '9' {
		return "m-" + id
	}
	return id
}
