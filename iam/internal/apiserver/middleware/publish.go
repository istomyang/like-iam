package middleware

import (
	"github.com/gin-gonic/gin"
	"istomyang.github.com/like-iam/component/pkg/middleware"
	"istomyang.github.com/like-iam/iam/internal/pkg"
	"istomyang.github.com/like-iam/log"
)

func NewSecretPublishInfo() middleware.PublishInfoInterface {
	return newPubInfo(pkg.PubSubChannel, pkg.MessageSecret)
}

func NewPolicyPublishInfo() middleware.PublishInfoInterface {
	return newPubInfo(pkg.PubSubChannel, pkg.MessagePolicy)
}

type pubInfo struct {
	c       *gin.Context
	channel string
	message string
}

func newPubInfo(channel, message string) *pubInfo {
	return &pubInfo{channel: channel, message: message}
}

func (i *pubInfo) WithContext(c *gin.Context) {
	i.c = c
}

func (i *pubInfo) Channel() string {
	return i.channel
}

func (i *pubInfo) Message() string {
	return i.message
}

func (i *pubInfo) Success() {
	log.L(i.c).Debugw("publish redis message", "method", i.c.Request.Method, "command", i.Message())
}

func (i *pubInfo) Fail(err error) {
	log.L(i.c).Errorw("publish redis message failed", "error", err.Error())
}
