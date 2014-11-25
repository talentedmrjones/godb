package server

// Command holds information about commands coming from client and facilitates communication between client and table
type Command struct {
	command string // create,read,update,delete
	data		[]byte // the data send by client, being acted upon by the command
	client  *Client // reference to the client so the Table can push to its reply channel
}

func (cmd *Command) Command () (string) {
	return cmd.command
}

func (cmd *Command) Data () ([]byte) {
	return cmd.data
}

func (cmd *Command) Client () (*Client) {
	return cmd.client
}
