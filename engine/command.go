package engine

type Command struct {
	Id			string
	Action 	string
	Db			string
	Table 	string
	Query		map[string][]byte
	client 	*Client
}
