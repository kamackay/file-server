package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.com/kamackay/go-api/logging"
)

type Authorizer struct {
	log *logrus.Logger
}

func New() *Authorizer {
	return &Authorizer{
		log: logging.GetLogger(),
	}
}

func (this *Authorizer) AllowedToDownload(ctx *gin.Context) bool {
	return true
}
