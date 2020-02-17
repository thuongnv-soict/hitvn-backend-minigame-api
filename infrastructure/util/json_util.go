package util

import (
	"encoding/json"
)

/**
 * Converts an object to string
 * @Param {object} v
 * @Return {string}
 */
func ToJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		// fmt.Println(err)
		return ""
	}

	return string(b)
}