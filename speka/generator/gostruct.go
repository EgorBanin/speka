package generator

import (
	"fmt"
	"io"
	"regexp"
	"speka/speka"
	"strings"
)

type goStruct struct {
	name   string
	fields []goStructField
}

type goStructField struct {
	name      string
	t         string
	jsonName  string
	validator string
}

type GoStruct struct {
	pckg  string
	types []goStruct
}

func NewGoStruct(pckg string) *GoStruct {
	return &GoStruct{
		pckg: pckg,
	}
}

func (g *GoStruct) Generate(p *speka.Property, w io.Writer) error {
	g.collectStructs(p)
	fmt.Fprintf(w, "package %s\n\n", g.pckg)
	for _, t := range g.types {
		fmt.Fprintf(w, "type %s struct {\n", t.name)
		for _, f := range t.fields {
			fmt.Fprintf(w, "\t%s %s `json:\"%s\"%s`\n", f.name, f.t, f.jsonName, f.validator)
		}
		fmt.Fprint(w, "}\n\n")
	}

	return nil
}

func (g *GoStruct) collectStructs(p *speka.Property) error {
	if p.Kind != speka.KindObject {
		return nil
	}

	st := goStruct{
		name:   camelCase(p.Name),
		fields: make([]goStructField, 0, len(p.Properties)),
	}

	for _, pp := range p.Properties {
		switch pp.Kind {
		case speka.KindObject:
			g.collectStructs(pp)
		case speka.KindArray:
			g.collectStructs(pp.Items)
		}

		var validator string
		rules := make([]string, 0)
		if pp.Required {
			rules = append(rules, "required")
		}
		if len(rules) > 0 {
			validator = fmt.Sprintf(" validate:\"%s\"", strings.Join(rules, ","))
		}

		st.fields = append(st.fields, goStructField{
			name:      camelCase(pp.Name),
			t:         getType(pp),
			jsonName:  pp.Name,
			validator: validator,
		})
	}

	g.types = append(g.types, st)

	return nil
}

var splitRegex = regexp.MustCompile("[^a-zA-Z]+")

func camelCase(s string) string {
	var result string
	ss := splitRegex.Split(s, -1)
	for i := range ss {
		if ss[i] == "" {
			continue
		}

		result += strings.ToUpper(ss[i][:1]) + ss[i][1:]
	}

	return result
}

func getType(p *speka.Property) string {
	t := "any"
	switch p.Kind {
	case speka.KindObject:
		t = camelCase(p.Name)
	case speka.KindString:
		t = "string"
	case speka.KindInteger:
		t = "int"
	case speka.KindNumber:
		t = "float64"
	case speka.KindArray:
		t = fmt.Sprintf("[]%s", getType(p.Items))
	}

	if !p.Required && p.Kind != speka.KindArray {
		t = "*" + t
	}

	return t
}
