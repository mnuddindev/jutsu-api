package parsers

import "encoding/json"

type StreamLink struct {
	File string `json:"file"`
	Type string `json:"type"`
}

type DecryptedSources struct {
	ID     string          `json:"id"`
	Type   string          `json:"type"`
	Link   StreamLink      `json:"link"`
	Tracks json.RawMessage `json:"tracks"`
	Intro  json.RawMessage `json:"intro"`
	Outro  json.RawMessage `json:"outro"`
	Iframe string          `json:"iframe"`
	Server string          `json:"server"`
}
