package engine

/*
Command is used to relay commands and data from connection to table channels
*/
type Command struct {
	ID         string                 `json:"id"`     // derived from client to be correlated to reply
	Action     string                 `json:"action"` // so far only "c" for create
	Db         string                 `json:"db"`
	Table      string                 `json:"table"`
	Data       map[string]interface{} `json:"data"`
	connection *Connection
}

/*
NewCommand ...
*/
func NewCommand(connection *Connection) *Command {
	command := &Command{}
	command.connection = connection
	return command
}
