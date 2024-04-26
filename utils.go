package api_client_go

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
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

func WriteBytesToFile(b []byte, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(b)
	return err
}

func WriteIndentedJSONToFile(data interface{}, path string) error {
	// Marshal the data with indentation
	jsonData, err := json.MarshalIndent(data, "", "  ")
	log.Infof("Writing Marshalled data: %s to file", jsonData)
	if err != nil {
		return err
	}

	return WriteBytesToFile(jsonData, path)
}

func ReadBytesFromFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func ReadJSONFromFile(path string, data interface{}) error {
	content, err := ReadBytesFromFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, data)
}
