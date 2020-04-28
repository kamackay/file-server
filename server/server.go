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
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

type Server struct {
	log        *logrus.Logger
	engine     *gin.Engine
	root       string
	store      *persistence.InMemoryStore
	cronRunner *cron.Cron
	auth       *auth.Authorizer
}

func New(root string) *Server {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	return &Server{
		log:        log,
		engine:     gin.Default(),
		root:       root,
		cronRunner: cron.New(),
		auth:       auth.New(root),
	}
}

func (this *Server) Start() {
	this.store = persistence.NewInMemoryStore(time.Second)
	this.engine.Use(this.auth.Bind())
	this.engine.Use(gzip.Gzip(gzip.BestCompression))
	this.engine.Use(cors.Default())

	this.engine.PUT("/*root", this.uploadFile())
	this.engine.POST("/*root", this.uploadFile())

	this.engine.GET("/*root", cache.CachePage(this.store, 5*time.Minute,
		func(ctx *gin.Context) {
			filename := this.root + ctx.Request.URL.Path
			if fi, exists, err := this.getFile(filename); err != nil {
				if exists {
					this.unknownError(ctx)
				} else {
					ctx.String(404, "File Not Found")
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
					if file, err := files.ReadFile(filename); err != nil {
						this.unknownError(ctx)
					} else {
						ctx.Data(200, file.ContentType, []byte(file.Data))
					}
				}
			}
		}))

	this.engine.DELETE("/*root", func(ctx *gin.Context) {
		file := this.root + ctx.Request.URL.Path
		if fi, exists, err := this.getFile(file); err != nil {
			if exists {
				ctx.String(404, "File Not Found")
			} else {
				this.unknownError(ctx)
			}
		} else {
			if fi.IsDir() {
				ctx.String(400, "Cannot Delete Folder")
			} else {
				if err := os.Remove(file); err != nil {
					ctx.String(500, "Could Not Delete File")
				} else {
					ctx.String(200, "Successfully Deleted File")
				}
			}
		}
	})

	if err := this.engine.Run(); err != nil {
		panic(err)
	} else {
		this.log.Info("Successfully Started Server")
	}
}

func (this *Server) uploadFile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		filename := this.root + ctx.Request.URL.Path

		if data, err := ctx.GetRawData(); err != nil {
			ctx.String(400, "Could Not Read File")
		} else {
			dir := path.Dir(filename)
			if err := os.MkdirAll(dir, 0644); err != nil {
				// Try anyways
				fmt.Println(err)
			}
			err = files.WriteFile(files.File{
				Data:        string(data),
				ContentType: ctx.ContentType(),
				LastUpdated: time.Now().UnixNano(),
				Name:        filename,
			})
			if err == nil {
				ctx.String(200, "Written Successfully")
				_ = this.store.Delete(cache.CreateKey(ctx.Request.RequestURI))
				_ = this.store.Delete(cache.CreateKey("/"))
			} else {
				this.log.Error("Error Writing File", err)
				ctx.String(500, "Error Writing File")
			}
		}
	}
}

func (this *Server) readConfig() *Server {
	var config Config
	if bytes, err := ioutil.ReadFile("/config.yml"); err != nil {
		this.log.Warnf("Could Not Read Config File")
		return this
	} else if err := yaml.Unmarshal(bytes, &config); err != nil {
		this.log.Warnf("Could not Unmarshal Config Object")
		return this
	}
	return this
}

type Config struct {
}
