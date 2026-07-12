package html

import (
	"embed"
	"fmt"
	"html/template"
	"io"

	"github.com/egorbanin/speka/speka"
)

//go:embed tpl/*.html
var tpl embed.FS

type Generator struct {
	tpl *template.Template
}

func NewGenerator(tpl *template.Template) *Generator {
	return &Generator{
		tpl: tpl,
	}
}

func CreateGenerator() (*Generator, error) {
	t, err := template.ParseFS(tpl, "tpl/*.html")
	if err != nil {
		return nil, err
	}

	return NewGenerator(t), nil
}

func (g *Generator) Html(s speka.Speka, w io.Writer) error {
	menu := make(menu, 0, len(s.Methods))
	i := 0
	for name := range s.Methods {
		menu = append(menu, menuItem{
			Text: name,
			Link: fmt.Sprintf("%d", i),
		})
		i++
	}

	g.tpl.ExecuteTemplate(w, "layout.html", layout{
		Menu: menu,
	})

	return nil
}
