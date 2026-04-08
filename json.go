package utils

import "encoding/json"

func JsonString(inter interface{}) string {
	interBytes, _ := json.Marshal(inter)
	return string(interBytes)
}
