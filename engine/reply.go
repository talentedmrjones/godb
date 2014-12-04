package engine

// Reply is sent back to client
type Reply struct {
	Id			string
	Status  uint16
	Result	Records
	Error		string
}

func NewReply (id string) *Reply {
	reply := &Reply{}
	reply.Id = id
	return reply
}
