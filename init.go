package sctl

// Init Struct for initializing a project
type Init struct {
	Project Project `json:"project"`
	Master  Node    `json:"master"`
}
