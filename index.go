package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	brotli "github.com/Solorad/gin-brotli"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"gitlab.com/kamackay/go-api/logging"
)

var root = os.Getenv("ROOT_PATH")

func main() {
	if len(root) == 0 {
		root = "/files"
	}
	log := logging.GetLogger()
	// Instantiate a new router
	engine := gin.Default()
	engine.Use(brotli.Brotli(brotli.Options{
		Quality: 7, // Default: 4
		LGWin:   11,
	}))
	engine.Use(gzip.Gzip(gzip.BestCompression))
	engine.Use(cors.Default())
	engine.Use(logger.SetLogger())

	engine.PUT("/*root", UploadFile)
	engine.POST("/*root", UploadFile)

	engine.GET("/*root", func(ctx *gin.Context) {
		file := root + ctx.Request.URL.Path
		if fi, exists, err := GetFile(file); err != nil {
			if exists {
				ctx.String(500, "Unknown Filesystem issue")
			} else {
				ctx.String(404, "File Not Found")
			}
		} else {
			if fi.IsDir() {
				if files, err := ioutil.ReadDir(file); err != nil {
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
				ctx.File(file)
			}
		}
	})

	engine.DELETE("/*root", func(ctx *gin.Context) {
		file := root + ctx.Request.URL.Path
		if fi, exists, err := GetFile(file); err != nil {
			if exists {
				ctx.String(404, "File Not Found")
			} else {
				ctx.String(500, "Unknown Filesystem issue")
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

	if err := engine.Run(); err != nil {
		panic(err)
	} else {
		log.Info("Successfully Started Server")
	}
}

func UploadFile(ctx *gin.Context) {
	filename := root + ctx.Request.URL.Path
	if data, err := ctx.GetRawData(); err != nil {
		ctx.String(400, "Could Not Read File")
	} else {
		dir := path.Dir(filename)
		if err := os.Mkdir(dir, 0644); err != nil {
			// Try anyways
			fmt.Println(err)
		}
		err = ioutil.WriteFile(filename, data, 0644)
		if err == nil {
			ctx.String(200, "Written Successfully")
		} else {
			ctx.String(500, "Error Writing File")
		}
	}
}
