package auth

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Authorizer struct {
	log    *logrus.Logger
	config *AuthConfig
	fsRoot string
}

func New(fsRoot string) *Authorizer {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	config := readConfigFile(log)
	return &Authorizer{
		log:    log,
		config: config,
		fsRoot: fsRoot,
	}
}

func (this *Authorizer) Bind() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fsPath := this.fsRoot + ctx.Request.URL.Path
		if fi, err := os.Stat(fsPath); err == nil && fi.IsDir() {
			// This is a directory
			ctx.Next()
		} else {
			if this.Validate(ctx) {
				ctx.Next()
			} else {
				return
			}
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
	return true
}

func (this *Authorizer) requiresValidation(ctx *gin.Context) bool {
	switch strings.ToUpper(ctx.Request.Method) {
	case "GET":
		return true
	default:
		return false
	}
}

func (this *Authorizer) validate(ctx *gin.Context) bool {
	authHeader := ctx.GetHeader("Authorization")
	validAuth := "Basic " +
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s",
			this.config.DefaultAuth.Username,
			this.config.DefaultAuth.Password)))
	return authHeader == validAuth
}

type AuthConfig struct {
	DefaultAuth struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"defaultCreds"`
}

func readConfigFile(log *logrus.Logger) *AuthConfig {
	var config AuthConfig
	if bytes, err := ioutil.ReadFile("/auth.yml"); err != nil {
		log.Warnf("Could Not Read Config File")
		return &config
	} else if err := yaml.Unmarshal(bytes, &config); err != nil {
		log.Warnf("Could not Unmarshal Config Object")
		return &config
	}
	return &config
}
