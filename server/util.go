package server

import (
	"github.com/gin-gonic/gin"
	"os"
)

func (this *Server) error(ctx *gin.Context, message string) {
	ctx.String(500, message)
}

func (this *Server) unknownError(ctx *gin.Context) {
	this.error(ctx, "Unknown Error")
}

func (this *Server) getFile(path string) (os.FileInfo, bool, error) {
	if fi, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fi, false, err
		} else {
			return fi, true, err
		}
	} else {
		return fi, true, err
	}
}
