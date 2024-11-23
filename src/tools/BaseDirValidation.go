package tools

import (
	"fmt"
	"path/filepath"
)

func IsWithinBaseDir(baseDir, targetPath string) (bool, error) {
	// Получаем абсолютные пути
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return false, err
	}
	absTargetPath, err := filepath.Abs(targetPath)
	if err != nil {
		return false, err
	}

	fmt.Println(absBaseDir)
	fmt.Println(absTargetPath)
	// Проверяем, что целевой путь начинается с базовой директории
	relPath, err := filepath.Rel(absBaseDir, absTargetPath)
	if err != nil {
		return false, err
	}
	if relPath == ".." || relPath[:3] == "../" {
		fmt.Println(relPath)
		return false, nil
	}
	return true, nil
}
