package client

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func FileExistsAndNotEmpty(filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return fileInfo.Size() != 0
}

func BuildDirectoryStructureIfNotExists(path string) error {
	if FileExists(path) {
		log.Infof("Directory %s already exists", path)
		return nil
	}
	return os.MkdirAll(path, os.ModePerm)
}
