package speka

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/hjson/hjson-go/v4"
)

const (
	KindObject  = "object"
	KindArray   = "array"
	KindString  = "string"
	KindNumber  = "number"
	KindInteger = "integer"

	FormatDate     = "date"
	FormatDateTime = "date-time"

	markNullable = "?"
)

type Property struct {
	Name       string
	Kind       string
	Enum       []string
	Format     string
	Properties []*Property
	Items      *Property
	Example    any
	Required   bool
}

func ParseProperty(name string, data any) (*Property, error) {
	p := new(Property)
	if strings.HasSuffix(name, markNullable) {
		p.Name = name[:len(name)-1]
	} else {
		p.Name = name
		p.Required = true
	}
	switch d := data.(type) {
	case *hjson.OrderedMap:
		p.Kind = KindObject
		p.Properties = make([]*Property, 0, d.Len())
		for _, n := range d.Keys {
			pp, err := ParseProperty(n, d.Map[n])
			if err != nil {
				return nil, err
			}

			p.Properties = append(p.Properties, pp)
		}
	case []any:
		p.Kind = KindArray
		if len(d) > 0 {
			items, err := ParseProperty(p.Name, d[0])
			if err != nil {
				return nil, err
			}

			p.Items = items
		}

	case string:
		p.Kind = KindString
		if strings.Contains(d, "|") {
			p.Enum = strings.Split(d, "|")
			p.Example = p.Enum[0]
		} else {
			_, err := time.Parse(time.DateOnly, d)
			if err == nil {
				p.Format = FormatDate
				p.Example = d
			}

			_, err = time.Parse(time.RFC3339, d)
			if err == nil {
				p.Format = FormatDateTime
				p.Example = d
			}
		}
	case float64:
		p.Kind = KindNumber
		if d == math.Trunc(d) {
			p.Kind = KindInteger
		}
	default:
		return nil, fmt.Errorf("can't parse %T", data)
	}

	return p, nil
}
