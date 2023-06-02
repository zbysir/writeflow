package model

type Category struct {
	Key  string  `json:"key"`
	Name Locales `json:"name"`
	Desc Locales `json:"desc"`
}
