package engine

type Command struct {
	Id			string
	Action 	string
	Db			string
	Table 	string
	Data		map[string][]byte
	client 	*Client
}
