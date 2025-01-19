package initializers

import (
	"github.com/gin-gonic/gin"
)

func SetupServer(server *gin.Engine) {
	server.SetTrustedProxies(nil)
}

