package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type Authorizer struct {
	log *logrus.Logger
}

func New() *Authorizer {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	return &Authorizer{
		log: log,
	}
}

func (this *Authorizer) Bind() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if this.Validate(ctx) {
			ctx.Next()
		} else {
			return
		}
	}
}

//Validate Return True if this request was valid
func (this *Authorizer) Validate(ctx *gin.Context) bool {
	this.log.Infof("Validating Request on `%s`", ctx.Request.URL.Path)
	if this.requiresValidation(ctx) {
		if this.validate(ctx) {
			return true
		} else {
			this.Decline(ctx)
			return false
		}
	} else {
		return true
	}
}

func (this *Authorizer) Decline(ctx *gin.Context) {
	ctx.Header("WWW-Authenticate", "Basic")
	ctx.AbortWithStatus(http.StatusUnauthorized)
}

func (this *Authorizer) AllowedToViewFolder(ctx *gin.Context) bool {
	return false
}

func (this *Authorizer) requiresValidation(ctx *gin.Context) bool {
	method := ctx.Request.Method
	return strings.ToUpper(method) != "GET"
}

func (this *Authorizer) validate(ctx *gin.Context) bool {
	authHeader := ctx.GetHeader("Authorization")
	return len(authHeader) > 0
}
