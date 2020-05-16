package server

import (
	"fmt"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gitlab.com/kamackay/filer/auth"
	"gitlab.com/kamackay/filer/files"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	CacheTime = 1 * time.Millisecond
)

type Server struct {
	log        *logrus.Logger
	engine     *gin.Engine
	root       string
	store      *persistence.InMemoryStore
	cronRunner *cron.Cron
	auth       *auth.Authorizer
	config     Config
}

func New(root string) *Server {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	return (&Server{
		log:        log,
		engine:     gin.Default(),
		root:       root,
		cronRunner: cron.New(),
		auth:       auth.New(root),
	}).readConfig()
}

func (this *Server) Start() {
	this.store = persistence.NewInMemoryStore(CacheTime)
	this.engine.Use(this.auth.Bind())
	this.engine.Use(gzip.Gzip(gzip.BestCompression))
	this.engine.Use(cors.Default())

	this.engine.PUT("/*root", this.uploadFile())
	this.engine.POST("/*root", this.postFile())

	this.engine.GET("/*root", cache.CachePage(this.store, CacheTime,
		func(ctx *gin.Context) {
			filename := this.root + ctx.Request.URL.Path
			urlPath := ctx.Request.URL.Path
			if regexp.MustCompile("^/ui/?.*").MatchString(urlPath) {
				if urlPath == "/ui" || urlPath == "/ui/" {
					// Send Root UI file
					this.sendFileNoMeta(ctx, "/ui/index.html", "text/html")
				} else {
					this.sendFileNoMeta(ctx, urlPath, "text/javascript")
				}
			} else if fi, exists, err := this.getFile(filename); err != nil {
				if exists {
					this.unknownError(ctx, err)
				} else {
					ctx.String(http.StatusNotFound, "File Not Found")
				}
			} else {
				if fi.IsDir() {
					if fs, err := ioutil.ReadDir(filename); err != nil {
						ctx.String(500, "Could Not get Contents of Directory")
					} else {
						if !this.auth.AllowedToViewFolder(ctx) {
							ctx.JSON(200, make([]string, 0))
							return
						}
						paths := make([]*files.JSONFile, 0)
						for _, f := range fs {
							if !strings.HasSuffix(f.Name(), files.MetaSuffix) {
								jsonData, err := files.GetJsonData(path.Join(filename, f.Name()))
								if err != nil {
									paths = append(paths, &files.JSONFile{
										Name:        f.Name(),
										ContentType: "text/plain",
										LastUpdated: 0,
									})
								} else {
									paths = append(paths, jsonData)
								}
							}
						}
						ctx.JSON(200, paths)
					}
				} else {
					this.sendFile(ctx, filename)
				}
			}
		}))

	this.engine.DELETE("/*root", func(ctx *gin.Context) {
		file := this.root + ctx.Request.URL.Path
		if fi, exists, err := this.getFile(file); err != nil {
			if exists {
				ctx.String(http.StatusNotFound, "File Not Found")
			} else {
				this.unknownError(ctx, err)
			}
		} else {
			if fi.IsDir() {
				ctx.String(http.StatusBadRequest, "Cannot Delete Folder")
			} else {
				if err := os.Remove(file); err != nil {
					ctx.String(500, "Could Not Delete File")
				} else {
					ctx.String(http.StatusOK, "Successfully Deleted File")
				}
			}
		}
	})

	if err := this.engine.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		panic(err)
	} else {
		this.log.Info("Successfully Started Server")
	}
}

func (this *Server) sendFile(ctx *gin.Context, filename string) {
	if file, reader, err := files.GetFile(filename); err != nil {
		this.unknownError(ctx, err)
	} else {
		defer reader.Close()
		fileSize := determineFileSize(file, reader)

		if fileSize < files.GetBufferLimit() {
			this.sendFileNoMeta(ctx, filename, file.ContentType)
			return
		}

		ctx.DataFromReader(http.StatusOK,
			fileSize,
			file.ContentType,
			reader,
			map[string]string{})

		// It seems that the Gin Cache is unable to handle this, and causes errors
		defer func() {
			_ = this.store.Delete(cache.CreateKey(ctx.Request.RequestURI))
		}()
	}
}

func (this *Server) sendFileNoMeta(ctx *gin.Context, filename string, contentType string) {
	data, err := ioutil.ReadFile(filename)
	if err == nil {
		ctx.Data(200, contentType, data)
	}
}

func determineFileSize(meta *files.MetaData, file *os.File) int64 {
	if meta.Size > 0 {
		return meta.Size
	} else {
		stat, err := file.Stat()
		if err != nil {
			return 0
		} else {
			return stat.Size()
		}
	}
}

// Use POST method to request server to download files from other location
func (this *Server) postFile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		switch ctx.ContentType() {
		case "text/plain":
			if data, err := ctx.GetRawData(); err != nil {
				ctx.String(500, "Unable to read URL")
			} else if err := files.DownloadFile(string(data), this.root+ctx.Request.URL.Path); err != nil {
				this.log.Error("Error Pulling MetaData", err)
				ctx.String(500, "Unable to download MetaData")
			} else {
				this.success(ctx)
			}
			return
		default:
			ctx.String(400, "Unsure how to use Request")
		}
	}
}

func (this *Server) success(ctx *gin.Context) {
	ctx.String(200, "Written Successfully")
	_ = this.store.Delete(cache.CreateKey(ctx.Request.RequestURI))
	_ = this.store.Delete(cache.CreateKey("/"))
}

func (this *Server) uploadFile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		filename := this.root + ctx.Request.URL.Path

		if strings.Contains(ctx.ContentType(), "multipart/form-data") {
			multi, err := ctx.Request.MultipartReader()
			if err != nil {
				this.unknownError(ctx, err)
				return
			}
			for {
				mimePart, err := multi.NextPart()
				if err == io.EOF {
					break
				}
				if err != nil {
					this.log.Warnf("Error reading multipart section: %v", err)
					this.unknownError(ctx, err)
					break
				}
				data, err := ioutil.ReadAll(mimePart)
				if err != nil {
					this.unknownError(ctx, err)
					this.log.Warnf("Unable to Read Error: %v", err)
					break
				}
				err = files.WriteFile(files.MetaData{
					ContentType: mime.TypeByExtension(filepath.Ext(mimePart.FileName())),
					LastUpdated: time.Now().UnixNano(),
					Name:        filename,
					Protected:   false,
				}, data)
				if err == nil {
					this.success(ctx)
				} else {
					this.log.Error("Error Writing MetaData", err)
					ctx.String(500, "Error Writing MetaData")
				}
			}
			return
		}
		if data, err := ctx.GetRawData(); err != nil {
			ctx.String(400, "Could Not Read MetaData")
		} else {
			dir := path.Dir(filename)
			if err := os.MkdirAll(dir, 0644); err != nil {
				// Try anyways
				fmt.Println(err)
			}
			err = files.WriteFile(files.MetaData{
				ContentType: ctx.ContentType(),
				LastUpdated: time.Now().UnixNano(),
				Name:        filename,
				Protected:   false,
			}, data)
			if err == nil {
				this.success(ctx)
			} else {
				this.log.Error("Error Writing MetaData", err)
				ctx.String(500, "Error Writing MetaData")
			}
		}
	}
}

func (this *Server) readConfig() *Server {
	var config Config
	if bytes, err := ioutil.ReadFile("/config.yml"); err != nil {
		this.log.Warnf("Could Not Read Config MetaData")
		return this
	} else if err := yaml.Unmarshal(bytes, &config); err != nil {
		this.log.Warnf("Could not Unmarshal Config Object")
		return this
	}
	this.log.Infof("Successfully Read Config File: %s", config)
	this.config = config
	return this
}

type Config struct {
	CacheServers []string `yaml:"cacheServers"`
}
