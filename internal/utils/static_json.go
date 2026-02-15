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

		// candidate paths (priority order)
		paths := []string{}

		// 1) container: next to binary
		if exec, err := os.Executable(); err == nil {
			base := filepath.Dir(exec)
			paths = append(paths, filepath.Join(base, "config", "docs", "roots.json"))
		}

		// 2) working directory (air / go run)
		if wd, err := os.Getwd(); err == nil {
			paths = append(paths, filepath.Join(wd, "config", "docs", "roots.json"))
		}

		// 3) fallback: project root guess (for IDE runs)
		paths = append(paths, "config/docs/roots.json")

		for _, path := range paths {
			if file, err := os.ReadFile(path); err == nil {
				json.Unmarshal(file, &rootDoc)
				return
			}
		}

		rootDoc = map[string]interface{}{
			"error": "documentation missing",
			"tried": paths,
		}
	})

	return rootDoc
}
