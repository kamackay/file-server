package auth

import (
	"encoding/base64"
	"fmt"
	"gitlab.com/kamackay/filer/files"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
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
		if this.validate(ctx, strings.ToUpper(ctx.Request.Method) == http.MethodGet) {
			return true
		} else {
			this.Decline(ctx)
			this.log.Warnf("Declining Request with Authorization `%s`", ctx.GetHeader("Authorization"))
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
	path := ctx.Request.URL.Path
	switch strings.ToUpper(ctx.Request.Method) {
	case http.MethodGet:
		if regexp.MustCompile("^/ui/?.*").MatchString(path) {
			// UI Paths should always return
			return false
		}
		meta, err := files.ReadMetaFile(this.fsRoot + ctx.Request.URL.Path)
		if err != nil {
			return true
		} else {
			return meta.Protected
		}
	case http.MethodPut, http.MethodPost, http.MethodDelete:
		return true
	default:
		return true
	}
}

func (this *Authorizer) validate(ctx *gin.Context, allowReadOnly bool) bool {
	authHeader := ctx.GetHeader("Authorization")
	validAuth := this.getAuthHeader()
	this.log.Debugf("Comparing %s to %s", authHeader, validAuth)
	if authHeader == validAuth {
		this.log.Infof("Successful auth from Admin user")
		return true
	} else {
		return allowReadOnly && this.validateReadOnly(authHeader)
	}

}

func (this *Authorizer) validateReadOnly(header string) bool {
	this.log.Infof("Validating Read-Only access")
	for _, auth := range this.config.ReadOnlyAuth {
		validAuth := encodeAuth(auth)
		this.log.Debugf("Comparing %s to %s", header, validAuth)
		if validAuth == header {
			this.log.Infof("Successful Auth from User %s", auth.Username)
			return true
		}
	}
	return false
}

func (this *Authorizer) getAuthHeader() string {
	return encodeAuth(this.config.DefaultAuth)
}

func encodeAuth(auth Auth) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s",
		auth.Username,
		auth.Password)))
}

type Config struct {
	DefaultAuth  Auth   `yaml:"defaultCreds"`
	ReadOnlyAuth []Auth `yaml:"readOnlyCreds"`
}

type Auth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
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
		config)
	return &config
}
