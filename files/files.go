package files

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	MetaSuffix         = ".meta"
	DefaultPermissions = 0644
	ProxyFolder        = "/temp"
)

func GetBufferLimit() int64 {
	s := os.Getenv("BUFFER_LIMIT")
	n, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return 50
	} else {
		return n
	}
}

func WriteFile(file MetaData, content []byte) error {
	if err := ioutil.WriteFile(file.Name, content, DefaultPermissions); err != nil {
		return err
	} else {
		file.Size = GetSize(file.Name)
		if err := writeMetaFile(file); err != nil {
			return err
		} else {
			return nil
		}
	}
}

/*	bool = handled -> Whether or not the sending of the file has been handled by this function
	MetaData = meta -> Metadata file on the requested URL
	os.File = file -> The File Requested on the Filesystem
	error = err -> Any Errors during file fetch
*/
func GetFile(ctx *gin.Context, filename string) (bool, *MetaData, *os.File, error) {
	if strings.HasSuffix(filename, MetaSuffix) {
		return false, nil, nil, errors.New("attempt to read meta file")
	}
	if file, err := ReadMetaFile(filename); err != nil {
		return false, nil, nil, err
	} else if len(file.ProxyPath) > 0 {
		// Handle Proxy Path
		return handleProxy(ctx, file)
	} else if reader, err := os.Open(filename); err != nil {
		return false, nil, nil, err
	} else {
		return false, file, reader, nil
	}
}

func CreateProxy(proxyUrl string, fileUrl string) error {
	_, meta, err := DownloadTempFile(proxyUrl, fileUrl)
	if err != nil {
		return err
	} else {
		return writeMetaFile(MetaData{
			Name:        meta.Name,
			ContentType: meta.ContentType,
			LastUpdated: meta.LastUpdated,
			Protected:   false,
			Size:        0,
			ProxyPath:   proxyUrl,
		})
	}
}

func GuessFileType(filename string) string {
	return mime.TypeByExtension(filepath.Ext(filename))
}

func WriteMetaFileFor(filename string) error {
	return writeMetaFile(MetaData{
		Name:        filename,
		ContentType: GuessFileType(filename),
		LastUpdated: time.Now().UnixNano(),
		Protected:   false,
		Size:        GetSize(filename),
		ProxyPath:   "",
	})
}

func writeMetaFile(file MetaData) error {
	return writeMetaFileTo(file, file.Name+MetaSuffix)
}

func writeMetaFileTo(file MetaData, path string) error {
	if data, err := yaml.Marshal(file); err != nil {
		return err
	} else if err := ioutil.WriteFile(path, data, DefaultPermissions); err != nil {
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
			Name:        MakeRelative(file.Name),
			ContentType: file.ContentType,
			LastUpdated: file.LastUpdated,
			Folder:      false,
			Protected:   file.Protected,
			Size:        file.Size,
			Count:       1,
		}, nil
	}
}

func handleProxy(ctx *gin.Context, metaFile *MetaData) (bool, *MetaData, *os.File, error) {
	urlPath := ctx.Request.URL.Path
	fmt.Printf("Handling Proxying on %s\n", urlPath)
	tempPath := ProxyFolder + urlPath
	if !FileExists(ProxyFolder) {
		_ = os.Mkdir(ProxyFolder, 0777)
	}
	if FileExists(tempPath) {
		// File has been previously downloaded
		if reader, err := os.Open(tempPath); err != nil {
			return false, nil, nil, err
		} else {
			return false, metaFile, reader, nil
		}
	} else {
		// Download File and write the stream to the response
		data, meta, err := DownloadTempFile(metaFile.ProxyPath, tempPath)
		if err != nil {
			return false, nil, nil, err
		} else {
			ctx.Data(200, meta.ContentType, data)
			return true, nil, nil, nil
		}
	}
}

func ReadMetaFile(filename string) (*MetaData, error) {
	if data, err := ioutil.ReadFile(filename + MetaSuffix); err != nil {
		return nil, err
	} else {
		var file MetaData
		if err := yaml.Unmarshal(data, &file); err != nil {
			return nil, err
		} else {
			//fmt.Printf("META: %+v", file)
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
				Size:        GetSize(filename),
			})
		}
	}
}

func DownloadTempFile(url string, filename string) ([]byte, *MetaData, error) {
	tempFilename := ProxyFolder + filename
	if resp, err := http.Get(url); err != nil {
		return nil, nil, err
	} else {
		defer resp.Body.Close()
		file, err := os.Create(tempFilename)
		if err != nil {
			return nil, nil, err
		}
		data, err := ioutil.ReadAll(resp.Body)
		meta := MetaData{
			Name:        tempFilename,
			ContentType: resp.Header.Get("Content-Type"),
			LastUpdated: time.Now().UnixNano(),
			Protected:   false,
			ProxyPath:   url,
			Size:        GetSize(tempFilename),
		}
		go func() {
			_, err = bufio.NewWriter(file).Write(data)
			if err != nil {
				fmt.Printf("Error Writing file Async %+v", err)
			} else {
				err := writeMetaFileTo(meta, "/files"+filename+MetaSuffix)
				if err != nil {
					fmt.Printf("Error Writing file Async %+v", err)
				}
			}
		}()
		return data, &meta, nil
	}
}

func MakeRelative(filename string) string {
	s := regexp.MustCompile("^/files").ReplaceAllString(filename, "")
	if regexp.MustCompile("/.*").MatchString(s) {
		return s
	} else {
		return "/" + s
	}
}

type MetaData struct {
	Name        string `yaml:"name"`
	ContentType string `yaml:"contentType"`
	LastUpdated int64  `yaml:"lastUpdated"`
	Protected   bool   `yaml:"protected"`
	Size        int64  `yaml:"size"`
	ProxyPath   string `yaml:"proxyPath"`
}

type JSONFile struct {
	Name        string `json:"name"`
	ContentType string `json:"contentType"`
	LastUpdated int64  `json:"lastUpdated"`
	Folder      bool   `json:"folder"`
	Protected   bool   `json:"protected"`
	Size        int64  `json:"size"`
	Count       int    `json:"count"`
}
