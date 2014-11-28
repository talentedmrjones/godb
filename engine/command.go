package engine

type Command struct {
	// TODO support unique ID for async replies
	Action 	string
	Db			string
	Table 	string
	Data		map[string][]byte
	Client 	*Client
}
