package engine

// Reply is sent back to client
type Reply struct {
	Id			string
	Status  uint16
	Records	[]map[string][]byte
	Error		string
}
