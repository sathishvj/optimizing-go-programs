package main

import (
	"encoding/json"
)

//easyjson:json
type JSONData struct {
	Data []string
}

func unmarshaljsonFn() {
	var j JSONData
	json.Unmarshal([]byte(`{"Data" : ["One", "Two", "Three"]} `), &j)
}

func easyjsonFn() {
	d := &JSONData{}
	d.UnmarshalJSON([]byte(`{"Data" : ["One", "Two", "Three"]} `))
}
