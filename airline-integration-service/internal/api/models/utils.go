package models

type ESResult[T any] struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source T `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
