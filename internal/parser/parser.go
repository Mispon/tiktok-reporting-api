package parser

import "encoding/json"

type Parser interface {
	Parse(responseRaw string) (*Response, error)
}

type parser struct{}

// New - constructor
func New() Parser {
	return &parser{}
}

// Parse - parses raw response from TikTok API
func (p parser) Parse(responseRaw string) (*Response, error) {
	var response Response
	if err := json.Unmarshal([]byte(responseRaw), &response); err != nil {
		return nil, err
	}
	return &response, nil
}
