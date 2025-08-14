package workdir

import (
	"fmt"
	"strings"
)

// you can use this library freely: "github.com/otiai10/copy"

type WorkDir struct {
	root  string
	files map[string]string
	dirs  map[string]bool
}

func InitEmptyWorkDir() *WorkDir {
	wd := &WorkDir{
		root:  ".",
		files: make(map[string]string),
		dirs:  make(map[string]bool),
	}
	return wd
}

func (wd *WorkDir) CreateFile(path string) error {
	wd.files[path] = ""
	return nil
}

func (wd *WorkDir) CreateDir(path string) error {
	wd.dirs[path] = true
	return nil
}

func (wd *WorkDir) WriteToFile(path string, content string) error {
	if _, exists := wd.files[path]; !exists {
		return fmt.Errorf("file not found: %s", path)
	}
	wd.files[path] = content
	return nil
}

func (wd *WorkDir) AppendToFile(path string, content string) error {
	if _, exists := wd.files[path]; !exists {
		return fmt.Errorf("file not found: %s", path)
	}
	wd.files[path] += content
	return nil
}

func (wd *WorkDir) CatFile(path string) (string, error) {
	content, ok := wd.files[path]
	if !ok {
		return "", fmt.Errorf("file not found: %s", path)
	}
	return content, nil
}

func (wd *WorkDir) ListFilesIn(path string) ([]string, error) {
	files := []string{}
	for file := range wd.files {
		if strings.HasPrefix(file, path) {
			files = append(files, file)
		}
	}
	return files, nil
}

func (wd *WorkDir) ListFilesRoot() []string {
	files := []string{}
	for file := range wd.files {
		files = append(files, file)
	}
	return files
}

func (wd *WorkDir) Clone() *WorkDir {
	newWd := &WorkDir{
		root:  wd.root,
		files: make(map[string]string),
		dirs:  make(map[string]bool),
	}

	// Copy all files
	for path, content := range wd.files {
		newWd.files[path] = content
	}

	// Copy all directories
	for path, exists := range wd.dirs {
		newWd.dirs[path] = exists
	}

	return newWd
}
