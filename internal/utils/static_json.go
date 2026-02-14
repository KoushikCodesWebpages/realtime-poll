package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	rootDoc  map[string]interface{}
	loadOnce sync.Once
)

func GetRootDoc() map[string]interface{} {

	loadOnce.Do(func() {

		// get location of THIS source file
		_, filename, _, _ := runtime.Caller(0)

		// project root = two folders up from utils
		base := filepath.Dir(filename)
		root := filepath.Join(base, "..", "..")

		path := filepath.Join(root, "config", "docs", "roots.json")

		file, err := os.ReadFile(path)
		if err != nil {
			rootDoc = map[string]interface{}{
				"error": "documentation missing",
				"path":  path,
			}
			return
		}

		json.Unmarshal(file, &rootDoc)
	})

	return rootDoc
}
