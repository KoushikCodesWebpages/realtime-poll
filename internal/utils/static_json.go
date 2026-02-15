package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

var (
	rootDoc  map[string]interface{}
	loadOnce sync.Once
)
func GetRootDoc() map[string]interface{} {

	loadOnce.Do(func() {

		execPath, err := os.Executable()
		if err != nil {
			rootDoc = map[string]interface{}{
				"error": "cannot find executable",
			}
			return
		}

		base := filepath.Dir(execPath)
		path := filepath.Join(base, "config", "docs", "roots.json")

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
