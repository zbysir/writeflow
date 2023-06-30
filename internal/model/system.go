package model

import "encoding/json"

type PluginSource struct {
	Url    string `json:"url"`
	Enable bool   `json:"enable"`
}

type Setting struct {
	Plugins []PluginSource `json:"plugins,omitempty"`
}

func (s Setting) Merge(a Setting) Setting {
	bs, _ := json.Marshal(a)
	_ = json.Unmarshal(bs, &s)
	return s
}
