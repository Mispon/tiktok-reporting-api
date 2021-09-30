package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// SendPOST process POST request
func SendPOST(url string, jsonData []byte) (map[string]interface{}, error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)

	return result, err
}
