package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mi0772/nuke-go/types"
)

func ListKeys(c *gin.Context) {
	database := getDatabase(c)

	var keys = database.Keys()
	var r = types.KeysResponse{Keys: keys}

	c.JSON(http.StatusOK, gin.H{
		"keys": r,
	})
}

func PartitionDetails(c *gin.Context) {
	database := getDatabase(c)

	result := database.DetailsPartitions()
	c.JSON(http.StatusOK, gin.H{
		"details_partitition": result,
	})
}
