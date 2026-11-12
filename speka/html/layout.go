package html

import "github.com/egorbanin/speka/speka"

type layout struct {
	Name    string
	Desc    string
	Menu    menu
	Methods []methodView
}

type methodView struct {
	ID       string
	Name     string
	Desc     string
	Params   []*speka.Property
	Response []*speka.Property
	Errors   []string
}
