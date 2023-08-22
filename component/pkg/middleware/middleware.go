package middleware

import (
	"github.com/gin-gonic/gin"
	gindump "github.com/tpkeeper/gin-dump"
)

var PresetMiddlewares = map[string]gin.HandlerFunc{
	"cors":     Cors(),
	"nocache":  NoCache(),
	"recovery": gin.Recovery(),
	"dump":     gindump.Dump(),
}

var MustMiddlewares = []gin.HandlerFunc{RequestId(), Secure(), Logger()}
