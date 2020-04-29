package files

import (
	"bufio"
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	MetaSuffix         = ".meta"
	DefaultPermissions = 0644
	BufferLimit = 60 * 1024 * 1024
)

func WriteFile(file MetaData, content []byte) error {
	if err := ioutil.WriteFile(file.Name, content, DefaultPermissions);
		err != nil {
		return err
	} else {
		file.Size = getSize(file.Name)
		if err := writeMetaFile(file);
			err != nil {
			return err
		} else {
			return nil
		}
	}
}

func GetFile(filename string) (*MetaData, *os.File, error) {
	if strings.HasSuffix(filename, MetaSuffix) {
		return nil, nil, errors.New("attempt to read meta file")
	}
	if file, err := ReadMetaFile(filename); err != nil {
		return nil, nil, err
	} else if reader, err := os.Open(filename); err != nil {
		return nil, nil, err
	} else {
		return file, reader, nil
	}
}

func writeMetaFile(file MetaData) error {
	if data, err := yaml.Marshal(MetaData{
		Name:        file.Name,
		ContentType: file.ContentType,
		LastUpdated: file.LastUpdated,
		Size:        file.Size,
	}); err != nil {
		return err
	} else if err := ioutil.WriteFile(file.Name+MetaSuffix, data, DefaultPermissions);
		err != nil {
		return err
	} else {
		return nil
	}
}

func GetJsonData(filename string) (*JSONFile, error) {
	if file, err := ReadMetaFile(filename);
		err != nil {
		return nil, err
	} else {
		return &JSONFile{
			Name:        file.Name,
			ContentType: file.ContentType,
			LastUpdated: file.LastUpdated,
		}, nil
	}
}

func ReadMetaFile(filename string) (*MetaData, error) {
	if data, err := ioutil.ReadFile(filename + MetaSuffix);
		err != nil {
		return nil, err
	} else {
		var file MetaData
		if err := yaml.Unmarshal(data, &file);
			err != nil {
			return nil, err
		} else {
			return &file, nil
		}
	}
}

func DownloadFile(url string, filename string) error {
	if resp, err := http.Get(url); err != nil {
		return err
	} else {
		defer resp.Body.Close()
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		_, err = bufio.NewWriter(file).ReadFrom(resp.Body)
		if err != nil {
			return err
		} else {
			return writeMetaFile(MetaData{
				Name:        filename,
				ContentType: resp.Header.Get("Content-Type"),
				LastUpdated: time.Now().UnixNano(),
				Protected:   false,
				Size:        getSize(filename),
			})
		}
	}
}

func getSize(filename string) int64 {
	if stat, err := os.Stat(filename); err != nil {
		return 0
	} else {
		return stat.Size()
	}
}

type MetaData struct {
	Name        string `yaml:"name"`
	ContentType string `yaml:"contentType"`
	LastUpdated int64  `yaml:"lastUpdated"`
	Protected   bool   `yaml:"protected"`
	Size        int64  `yaml:"size"`
}

type JSONFile struct {
	Name        string `json:"name"`
	ContentType string `json:"contentType"`
	LastUpdated int64  `json:"lastUpdated"`
}
