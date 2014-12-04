package engine

type Record struct {
	Data	map[string][]byte
}

type Records []Record

func NewRecord () Record {
	return Record{make(map[string][]byte)}
}
