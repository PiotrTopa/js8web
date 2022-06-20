package main

type Js8callEvent struct {
	Type   string                 `json:"type"`
	Value  string                 `json:"value"`
	Params map[string]interface{} `json:"params"`
}
