// Code generated by "nestcsv"; DO NOT EDIT.

package table

type RewardParamValue struct {
	Str   string  `json:"Str"`
	Int   int     `json:"Int"`
	Float float64 `json:"Float"`
}

type Reward struct {
	Type       string           `json:"Type"`
	ParamValue RewardParamValue `json:"ParamValue"`
	ParamType  string           `json:"ParamType"`
}
