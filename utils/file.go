package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CreateFile(name string) (*os.File, error) {
	relPath := strings.TrimPrefix(name, filepath.Dir(name))
	relPath = strings.Replace(relPath, `\`, `/`, -1)
	relPath = strings.TrimLeft(relPath, `/`)
	relPath = strings.Replace(relPath, " ", "-", -1)
	return os.Create(fmt.Sprintf("%s.csv", relPath))
}
