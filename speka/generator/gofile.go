package generator

import (
	"fmt"
	"io"
	"regexp"
	"slices"
	"strings"

	"github.com/egorbanin/speka/speka"
)

const (
	typeStruct         = "struct"
	typeSlice          = "[]"
	typeJsonRawMessage = "json.RawMessage"
)

type GoFile struct {
	name             string
	generatedComment string
	packageName      string
	imports          []string
	types            []*goStruct
}

func NewGoFile(name, generatedComment, packageName string) *GoFile {
	return &GoFile{
		name:             name,
		generatedComment: generatedComment,
		packageName:      packageName,
		imports:          make([]string, 0),
		types:            make([]*goStruct, 0),
	}
}

func (f *GoFile) AddProperty(p *speka.Property, opts GoStructOpts) error {
	s, err := f.structs(p, nil, opts)
	if err != nil {
		return err
	}

	f.types = append(f.types, s...)

	return nil
}

func (f *GoFile) Write(w io.Writer) {
	fmt.Fprintf(w, "// %s\npackage %s\n\n", f.generatedComment, f.packageName)

	if len(f.imports) > 0 {
		fmt.Fprintf(w, "import (\n\t\"%s\"\n)\n\n", strings.Join(f.imports, "\n\t"))
	}

	for _, t := range f.types {
		tt := t.t
		if strings.HasPrefix(tt, "*") {
			tt = t.t[1:]
		}

		fmt.Fprintf(w, "type %s %s", t.name, tt)

		if tt != typeStruct {
			fmt.Fprintln(w)
			continue
		}

		fmt.Fprintln(w, " {")

		for _, f := range t.fields {
			fmt.Fprintf(w, "\t%s %s `json:\"%s\"%s`\n", f.name, f.t, f.jsonName, f.validator)
		}

		fmt.Fprintln(w, "}", "")
	}
}

func (f *GoFile) addImport(imports ...string) {
	for _, imp := range imports {
		if slices.Index(f.imports, imp) == -1 {
			f.imports = append(f.imports, imp)
		}
	}
}

func (f *GoFile) structs(p *speka.Property, parent *goStruct, opts GoStructOpts) ([]*goStruct, error) {
	if p == nil {
		return nil, nil
	}

	if p.Kind != speka.KindArray && p.Kind != speka.KindObject {
		return nil, nil
	}

	name := camelCase(p.Name)
	if parent != nil {
		name = parent.name + name
	}

	t := getType(p)
	switch t {
	case typeSlice:
		if parent != nil {
			t = t + parent.name
		}

		t = t + camelCase(p.Name)

		if p.Items != nil {
			t += camelCase(p.Items.Name)
		} else {
			t += "any"
		}

	case typeJsonRawMessage:
		if parent != nil {
			f.addImport("encoding/json")

			return nil, nil
		}

		t = "any"
	}

	current := &goStruct{
		name:   name,
		t:      t,
		fields: make([]goStructField, 0, len(p.Properties)),
	}

	for _, pp := range p.Properties {
		validator := ""
		v := make([]string, 0)
		if opts.Validator {
			if pp.Required {
				v = append(v, "required")
			} else {
				v = append(v, "omitempty")
			}

			if len(pp.Enum) > 0 {
				v = append(v, fmt.Sprintf("oneof=%s", strings.Join(pp.Enum, " ")))
			}

			if pp.Kind == speka.KindArray {
				v = append(v, "dive", "required")
			}

			switch pp.Format {
			case speka.FormatDate:
				v = append(v, "datetime=2006-01-02")
			case speka.FormatDateTime:
				v = append(v, "datetime=2006-01-02T15:04:05Z07:00")
			case speka.FormatEmail:
				v = append(v, "email")
			}
		}

		if len(v) > 0 {
			validator = fmt.Sprintf(" validate:\"%s\"", strings.Join(v, ","))
		}

		t := getType(pp)
		asterisk := ""
		if strings.HasPrefix(t, "*") {
			t = t[1:]
			asterisk = "*"
		}
		switch t {
		case typeStruct:
			t = asterisk + current.name + camelCase(pp.Name)
		case typeSlice:
			tt := getType(pp.Items)
			if tt == typeStruct {
				t += current.name + camelCase(pp.Name) + camelCase(pp.Items.Name)
			} else {
				t += tt
			}
		}

		current.fields = append(current.fields, goStructField{
			name:      camelCase(pp.Name),
			t:         t,
			jsonName:  pp.Name,
			validator: validator,
		})
	}

	s := []*goStruct{
		current,
	}

	if p.Kind == speka.KindArray {
		if parent != nil {
			s = []*goStruct{}
		}

		ss, err := f.structs(p.Items, current, opts)
		if err != nil {
			return nil, err
		}

		s = append(s, ss...)
	}

	for _, pp := range p.Properties {
		ss, err := f.structs(pp, current, opts)
		if err != nil {
			return nil, err
		}

		s = append(s, ss...)
	}

	return s, nil
}

type goStruct struct {
	name   string
	t      string
	fields []goStructField
}

type goStructField struct {
	name      string
	t         string
	jsonName  string
	validator string
}

type GoStructOpts struct {
	Validator bool
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
	if p == nil {
		return t
	}

	switch p.Kind {
	case speka.KindObject:
		if len(p.Properties) == 0 {
			t = typeJsonRawMessage
		} else {
			t = typeStruct
		}
	case speka.KindString:
		t = "string"
	case speka.KindInteger:
		t = "int"
	case speka.KindNumber:
		t = "float64"
	case speka.KindBoolean:
		t = "bool"
	case speka.KindArray:
		t = typeSlice
	}

	if !p.Required && t != typeSlice && t != typeJsonRawMessage {
		t = "*" + t
	}

	return t
}
