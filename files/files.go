package files

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

const (
	MetaSuffix         = ".meta"
	DefaultPermissions = 0644
)

func WriteFile(file File) error {
	if err := writeMetaFile(file); err != nil {
		return err
	} else if err := ioutil.WriteFile(file.Name, []byte(file.Data), DefaultPermissions); err != nil {
		return err
	} else {
		return nil
	}

}

func ReadFile(filename string) (*File, error) {
	if strings.HasSuffix(filename, MetaSuffix) {
		return nil, errors.New("attempt to read meta file")
	}
	if data, err := ioutil.ReadFile(filename); err != nil {
		return nil, err
	} else if file, err := ReadMetaFile(filename); err != nil {
		return nil, err
	} else {
		file.Data = string(data)
		return file, nil
	}
}

func writeMetaFile(file File) error {
	if data, err := yaml.Marshal(File{
		Name:        file.Name,
		Data:        "",
		ContentType: file.ContentType,
		LastUpdated: file.LastUpdated,
	}); err != nil {
		return err
	} else if err := ioutil.WriteFile(file.Name+MetaSuffix, data, DefaultPermissions); err != nil {
		return err
	} else {
		return nil
	}
}

func GetJsonData(filename string) (*JSONFile, error) {
	if file, err := ReadMetaFile(filename); err != nil {
		return nil, err
	} else {
		return &JSONFile{
			Name:        file.Name,
			ContentType: file.ContentType,
			LastUpdated: file.LastUpdated,
		}, nil
	}
}

func ReadMetaFile(filename string) (*File, error) {
	if data, err := ioutil.ReadFile(filename + MetaSuffix); err != nil {
		return nil, err
	} else {
		var file File
		if err := yaml.Unmarshal(data, &file); err != nil {
			return nil, err
		} else {
			return &file, nil
		}
	}
}

type File struct {
	Name        string `yaml:"name"`
	Data        string `yaml:"data"`
	ContentType string `yaml:"contentType"`
	LastUpdated int64  `yaml:"lastUpdated"`
}

type JSONFile struct {
	Name        string `json:"name"`
	ContentType string `json:"contentType"`
	LastUpdated int64  `json:"lastUpdated"`
}
