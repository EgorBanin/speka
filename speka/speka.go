package speka

import "github.com/hjson/hjson-go/v4"

type Speka struct {
	Name    string            `json:"name"`
	Methods map[string]Method `json:"methods"`
}

type Method struct {
	Rq *hjson.OrderedMap `json:"rq"`
	Rs *hjson.OrderedMap `json:"rs"`
}
