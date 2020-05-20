package compresssion

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.com/kamackay/filer/utils"
	"regexp"
)

const (
	Brotli = "br"
)

type Compressor struct {
	log *logrus.Logger
}

func New() *Compressor {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	return &Compressor{
		log: log,
	}
}

func (this *Compressor) Bind() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		switch this.getBestCompression(ctx) {
		case Brotli:
			// TODO Do compression
			return
		case "gz":
			// Already Have a Plugin to handle this
			return
		default:
			return
		}
	}
}

func (this *Compressor) getBestCompression(ctx *gin.Context) string {
	types := this.getAllowedCompressions(ctx)
	if utils.StrArrIncludes(types, Brotli) >= 0 {
		// Includes Brotli
		return Brotli
	}

	return ""
}

func (this *Compressor) getAllowedCompressions(ctx *gin.Context) []string {
	return regexp.MustCompile(",\\s*").
		Split(ctx.GetHeader("Accept-Encoding"), -1)
}