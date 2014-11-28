package engine

type Record struct {
	Db			string
	Table 	string
	Data		map[string][]byte
	Action 	string
	Client 	*Client
}
