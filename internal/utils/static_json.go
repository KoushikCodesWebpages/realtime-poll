package utils

import (
	"encoding/json"
	"os"
	"sync"
)

var (
	rootDoc  map[string]interface{}
	loadOnce sync.Once
)

func GetRootDoc() map[string]interface{} {

	loadOnce.Do(func() {
		file, err := os.ReadFile("configs/docs/root.json")
		if err != nil {
			rootDoc = map[string]interface{}{
				"error": "documentation missing",
			}
			return
		}

		json.Unmarshal(file, &rootDoc)
	})

	return rootDoc
}
