package common

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
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
	if err != nil {
		return err
	}

	log.Infof("Writing Marshalled data: %s to file", jsonData)

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

// GetCleanMapFromStruct returns a map with only the non-nil fields
// warning: this function will cause problems if trying to use default values in the struct
func GetCleanMapFromStruct(st interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	val := reflect.ValueOf(st).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		jsonTag := typeField.Tag.Get("json")

		if strings.Contains(jsonTag, ",") {
			jsonTagParts := strings.Split(jsonTag, ",")
			jsonTag = jsonTagParts[0]
		}

		if !valueField.IsZero() {
			if typeField.Name == "First" || typeField.Name == "Last" {
				m[jsonTag] = strconv.Itoa(int(valueField.Int()))
			} else {
				m[jsonTag] = valueField.String()
			}
		}
	}
	return m
}

func CloseGenericResponseWithLog(resp *http.Response, respErr error) {
	// Nothing to close
	if resp == nil {
		return
	}

	// Log if err is not nil
	if respErr != nil {
		log.Warnf("Error present together with response: %s", respErr)
	}

	method := ""
	endpoint := ""

	if resp.Request != nil {
		method = resp.Request.Method
		endpoint = resp.Request.URL.String()
	}

	log.Debugf("Closing response [%s]-%s", method, endpoint)
	err := resp.Body.Close()
	if err != nil {
		log.Warnf("Error closing responses %s body: %s", endpoint, err)
	}
}

func IsNilValueInterface(i interface{}) bool {
	iv := reflect.ValueOf(i)
	if !iv.IsValid() {
		return true
	}
	switch iv.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface:
		return iv.IsNil()
	default:
		return false
	}
}
