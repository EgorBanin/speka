package speka

import "github.com/hjson/hjson-go/v4"

type Speka struct {
	Name    string            `json:"name"`
	Methods map[string]Method `json:"methods"`
}

type Method struct {
	Rq *hjson.Node `json:"rq"`
	Rs *hjson.Node `json:"rs"`
}
