package servers

import (
	"github.com/gin-gonic/gin"
	"github.com/mi0772/nuke-go/engine"
	"github.com/mi0772/nuke-go/handlers"
	"github.com/mi0772/nuke-go/middleware"
)

func StartHTTPServer(database *engine.Database) {
	r := gin.Default()

	r.Use(middleware.DatabaseMiddleware(database))
	r.GET("/admin/keys", handlers.ListKeys)
	r.DELETE("/admin/clear", handlers.Clear)
	r.GET("/admin/partitions/details", handlers.PartitionDetails)
	r.POST("/push_file", handlers.PushFile)
	r.POST("/push_string", handlers.PushString)
	r.GET("/pop/:key", handlers.Pop)
	r.GET("/read/:key", handlers.Read)
	r.Run()
}
