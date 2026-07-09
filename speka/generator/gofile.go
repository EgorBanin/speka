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
	typeJsonRawMessage = "json.RawMessage"
)

type GoFile struct {
	name        string
	packageName string
	imports     []string
	types       []*goStruct
}

func NewGoFile(name, packageName string) *GoFile {
	return &GoFile{
		name:        name,
		packageName: packageName,
		imports:     make([]string, 0),
		types:       make([]*goStruct, 0),
	}
}

func (f *GoFile) AddProperty(p *speka.Property, opts GoStructOpts) error {
	s, err := f.structs(p, opts)
	if err != nil {
		return err
	}

	f.types = append(f.types, s...)

	return nil
}

func (f *GoFile) Write(w io.Writer) {
	fmt.Fprintf(w, "package %s\n\n", f.packageName)

	if len(f.imports) > 0 {
		fmt.Fprintf(w, "import (\n\t%s\n)\n\n", strings.Join(f.imports, "\n\t"))
	}

	for _, t := range f.types {
		fmt.Fprintf(w, "type %s %s", t.name, t.t)

		if t.t != typeStruct {
			fmt.Fprintln(w)
			continue
		}

		fmt.Fprintln(w, "{")

		for _, f := range t.fields {
			fmt.Fprintf(w, "\t%s %s\n", f.name, f.t)
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

func (f *GoFile) structs(p *speka.Property, opts GoStructOpts) ([]*goStruct, error) {
	if p == nil {
		return nil, nil
	}

	if p.Kind != speka.KindArray && p.Kind != speka.KindObject {
		return nil, nil
	}

	if p.Kind == speka.KindArray {
		s, err := f.structs(p.Items, opts)
		if err != nil {
			return nil, err
		}

		return s, nil
	}

	s := make([]*goStruct, 0, len(p.Properties))
	for _, pp := range p.Properties {
		ss, err := f.structs(pp, opts)
		if err != nil {
			return nil, err
		}

		s = append(s, ss...)
	}
	fields := make([]goStructField, 0, len(p.Properties))
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
			}
		}

		if len(v) > 0 {
			validator = fmt.Sprintf(" validate:\"%s\"", strings.Join(v, ","))
		}

		t := getType(pp)
		if t == typeStruct {
			t = camelCase(pp.Name)
		}

		fields = append(fields, goStructField{
			name:      camelCase(pp.Name),
			t:         t,
			jsonName:  pp.Name,
			validator: validator,
		})
	}

	t := getType(p)
	if t == typeJsonRawMessage {
		f.addImport("encoding/json")

		return nil, nil
	}

	s = append(s, &goStruct{
		name:   camelCase(p.Name),
		t:      t,
		fields: fields,
	})

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
		tt := getType(p.Items)
		if tt == typeStruct {
			t = fmt.Sprintf("[]%s", camelCase(p.Items.Name))
		} else {
			t = fmt.Sprintf("[]%s", tt)
		}
	}

	if !p.Required && p.Kind != speka.KindArray && t != typeJsonRawMessage {
		t = "*" + t
	}

	return t
}
