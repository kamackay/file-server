package server

import (
	"fmt"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gitlab.com/kamackay/filer/auth"
	"gitlab.com/kamackay/go-api/logging"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
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
	return &Server{
		log:        logging.GetLogger(),
		engine:     gin.Default(),
		root:       root,
		cronRunner: cron.New(),
		auth:       auth.New(),
	}
}

func (this *Server) Start() {
	this.store = persistence.NewInMemoryStore(time.Second)
	this.engine.Use(gzip.Gzip(gzip.BestCompression))
	this.engine.Use(cors.Default())
	this.engine.Use(logger.SetLogger())

	this.engine.PUT("/*root", this.uploadFile())
	this.engine.POST("/*root", this.uploadFile())

	this.engine.GET("/*root", cache.CachePage(this.store, time.Minute*5,
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
					if files, err := ioutil.ReadDir(filename); err != nil {
						ctx.String(500, "Could Not get Contents of Directory")
					} else {
						paths := make([]string, 0)
						for _, f := range files {
							paths = append(paths, f.Name())
						}
						fmt.Printf("Found %d paths", len(paths))
						ctx.JSON(200, paths)
					}
				} else {
					if data, err := ioutil.ReadFile(filename); err != nil {
						this.unknownError(ctx)
					} else {
						var file File
						if err := yaml.Unmarshal(data, &file); err != nil {
							this.unknownError(ctx)
						} else {
							ctx.Data(200, file.ContentType, []byte(file.Data))
						}
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
			file, err := yaml.Marshal(File{
				Data:        string(data),
				ContentType: ctx.ContentType(),
				LastUpdated: time.Now().UnixNano(),
			})
			if err != nil {
				fmt.Printf("Error Parsing into File Struct: %s", err)
				ctx.String(500, "Error Writing File")
				return
			}
			err = ioutil.WriteFile(filename, file, 0644)
			if err == nil {
				ctx.String(200, "Written Successfully")
			} else {
				ctx.String(500, "Error Writing File")
			}
		}
	}
}

type File struct {
	Data        string `yaml:"data"`
	ContentType string `yaml:"contentType"`
	LastUpdated int64  `yaml:"lastUpdated"`
}
