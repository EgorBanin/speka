package speka

import "github.com/hjson/hjson-go/v4"

type Speka struct {
	Name    string            `json:"name"`
	Desc    string            `json:"desc"`
	Methods map[string]Method `json:"methods"`
}

type Method struct {
	Desc   string      `json:"desc"`
	Rq     *hjson.Node `json:"rq"`
	Rs     *hjson.Node `json:"rs"`
	Errors []string    `json:"errors"`
}
