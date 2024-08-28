package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mi0772/nuke-go/engine"
	handlers2 "github.com/mi0772/nuke-go/servers/http/handlers"
	"github.com/mi0772/nuke-go/servers/http/middleware"
	"github.com/mi0772/nuke-go/types"
)

func StartHTTPServer(database *engine.Database, config *types.Configuration) {
	r := gin.Default()

	r.Use(middleware.DatabaseMiddleware(database))
	r.GET("/admin/keys", handlers2.ListKeys)
	r.DELETE("/admin/clear", handlers2.Clear)
	r.GET("/admin/partitions/details", handlers2.PartitionDetails)
	r.POST("/push_file", handlers2.PushFile)
	r.POST("/push_string", handlers2.PushString)
	r.GET("/pop/:key", handlers2.Pop)
	r.GET("/read/:key", handlers2.Read)
	r.Run(fmt.Sprintf(":%s", config.HttpServerPort))
}
