package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"strings"
)

type Authorizer struct {
	log    *logrus.Logger
	config *AuthConfig
}

func New() *Authorizer {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	config := readConfigFile(log)
	return &Authorizer{
		log:    log,
		config: config,
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

type AuthConfig struct {
	DefaultAuth string
}

func readConfigFile(log *logrus.Logger) *AuthConfig {
	var config AuthConfig
	if bytes, err := ioutil.ReadFile("/config.yml"); err != nil {
		log.Warnf("Could Not Read Config File")
		return &config
	} else if err := yaml.Unmarshal(bytes, &config); err != nil {
		log.Warnf("Could not Unmarshal Config Object")
		return &config
	}
	return &config
}
