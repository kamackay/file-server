package auth

import (
	"encoding/base64"
	"fmt"
	"gitlab.com/kamackay/filer/files"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Authorizer struct {
	log    *logrus.Logger
	config *Config
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
		this.log.Infof("Request Requires Validation")
		if this.validate(ctx) {
			return true
		} else {
			this.Decline(ctx)
			this.log.Warnf("Declining Request")
			return false
		}
	} else {
		this.log.Infof("Request Does Not Requires Validation")
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
		meta, err := files.ReadMetaFile(this.fsRoot + ctx.Request.URL.Path)
		if err != nil {
			return true
		} else {
			return meta.Protected
		}
	case "PUT", "POST", "DELETE":
		return true
	}
	return false
}

func (this *Authorizer) validate(ctx *gin.Context) bool {
	authHeader := ctx.GetHeader("Authorization")
	validAuth := "Basic " + this.getAuthHeader()
	this.log.Debugf("Comparing %s to %s", authHeader, validAuth)
	return authHeader == validAuth
}

func (this *Authorizer) getAuthHeader() string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s",
		this.config.DefaultAuth.Username,
		this.config.DefaultAuth.Password)))
}

type Config struct {
	DefaultAuth struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"defaultCreds"`
}

func readConfigFile(log *logrus.Logger) *Config {
	var config Config
	if bytes, err := ioutil.ReadFile("/auth.yml"); err != nil {
		log.Warnf("Could Not Read Config MetaData")
		return &config
	} else if err := yaml.Unmarshal(bytes, &config); err != nil {
		log.Warnf("Could not Unmarshal Config Object")
		return &config
	}
	log.Infof("Read Auth Config From Filesystem: %s",
		config.DefaultAuth)
	return &config
}
