package util

import "encoding/json"

func JSONstringify(data interface{}) string {
	if data == nil {
		return ""
	}
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return ""
	}
	return string(jsonData)
}
