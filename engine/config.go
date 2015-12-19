package engine

import (
	"encoding/json"
	"os"
)

// LoadConfig ...
func LoadConfig(path string) (JSON, error) {

	file, fileOpenErr := os.OpenFile(path, os.O_RDONLY, 0666)

	if fileOpenErr != nil {
		return nil, fileOpenErr
	}

	stats, statsError := file.Stat()
	if statsError != nil {
		return nil, statsError
	}

	bytes := make([]byte, stats.Size())
	_, readError := file.Read(bytes)
	if readError != nil {
		return nil, readError
	}

	config := make(JSON)

	bytesUnmarshalErr := json.Unmarshal(bytes, &config)
	if bytesUnmarshalErr != nil {
		return nil, bytesUnmarshalErr
	}

	return config, nil

}
