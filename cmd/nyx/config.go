package main

type configuration struct {
	Version string               `json:"version"`
	Checks  []checkConfiguration `json:"checks"`
}

type checkConfiguration struct {
	Name       string            `json:"name"`
	Check      string            `json:"check"`
	Help       string            `json:"help"`
	Identifier string            `json:"identifier"`
	Parameter  map[string]string `json:"parameter"`
}
