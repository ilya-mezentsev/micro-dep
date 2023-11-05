package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func MustParse[T any](path string) T {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("unable to read settings file: %v", err))
	}

	var s T
	err = json.Unmarshal(data, &s)
	if err != nil {
		panic(fmt.Errorf("unable to unmarshal settings to struct: %v", err))
	}

	return s
}
