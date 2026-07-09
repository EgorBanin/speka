package generator

import (
	"testing"

	"github.com/egorbanin/speka/speka"
	"github.com/stretchr/testify/assert"
)

func TestGoFile_AddProperty(t *testing.T) {
	p, _ := speka.ParseProperty("X", []any{
		map[string]any{
			"y": 1,
		},
	})

	f := NewGoFile("test", "test")
	err := f.AddProperty(p, GoStructOpts{})

	assert.NoError(t, err)
	assert.Equal(t, f.types, []*goStruct{
		{
			name: "X",
			t:    "[]XItem",
			fields: []goStructField{
				{
					name:     "Y",
					t:        "int",
					jsonName: "y",
				},
			},
		},
	})
}
