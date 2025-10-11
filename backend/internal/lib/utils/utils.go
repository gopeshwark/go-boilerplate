package utils

import (
	"encoding/json"
	"fmt"
)

func PringJSON(v interface{}) {
	json, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("ERror marshalling to JSON: ", err)
		return
	}
	fmt.Println("JSON: ", string(json))
}
