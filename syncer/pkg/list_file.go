package pkg

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

type fileInfo struct {
	path    string
	size    int64
	modTime time.Time
}

func ListFiles(suffixes []string, basePath string) ([]fileInfo, error) {
	infos, err := ioutil.ReadDir(basePath)
	if err != nil {
		return nil, err
	}
	var files []fileInfo
	for _, info := range infos {
		if info.Size() > 0 && matchSuffix(suffixes, info.Name()) {
			files = append(files, fileInfo{path: filepath.Join(basePath, info.Name()), size: info.Size(), modTime: info.ModTime()})
		}
	}
	return files, nil
}

func matchSuffix(suffixes []string, fileName string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(fileName, suffix) {
			return true
		}
	}
	return false
}
