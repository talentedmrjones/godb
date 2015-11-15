package engine

// Reply is sent back to client
type Reply struct {
	ID     string                 `json:"id"`
	Status uint16                 `json:"status"`
	Result map[string]interface{} `json:"result"`
	Error  string                 `json:"error"`
}

/*
NewReply creates a reply. The ID is provided by the command to link commands with the respective reply.
*/
func NewReply(id string) Reply {
	reply := Reply{}
	reply.ID = id
	return reply
}
