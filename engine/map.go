package engine

// Map ...
type Map struct {
	Type     string   `json:"type"`
	Interval float64  `json:"interval"`
	Fields   []string `json:"fields"`
}
