package config

import (
	"os"
	"path/filepath"
)

// FindModuleRoot ищет корень модуля (директорию с go.mod), поднимаясь от текущей рабочей директории вверх.
// Используется интеграционными тестами для загрузки .env из корня проекта.
func FindModuleRoot() (string, bool) {
	dir, err := os.Getwd()
	if err != nil {
		return "", false
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}
