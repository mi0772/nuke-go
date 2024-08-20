package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mi0772/nuke-go/engine"
	"github.com/mi0772/nuke-go/types"
)

func Pop(c *gin.Context) {
	database := getDatabase(c)
	key := c.Param("key")

	item, err := database.Pop(key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func Read(c *gin.Context) {
	database := getDatabase(c)
	key := c.Param("key")

	item, err := database.Read(key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func PushString(c *gin.Context) {
	database := getDatabase(c)
	var request types.PushItemStringRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.Push(request.Key, []byte(request.Value))
	c.JSON(http.StatusCreated, gin.H{})
}

func PushFile(c *gin.Context) {
	database := getDatabase(c)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer f.Close()

	fileBytes, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	metadataJSON := c.PostForm("request")
	var request types.PushItemFileRequest
	if err := json.Unmarshal([]byte(metadataJSON), &request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON in request"})
		return
	}

	i, err := database.Push(request.Key, fileBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	json, err := json.Marshal(i)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"item": json,
	})
}

func getDatabase(c *gin.Context) *engine.Database {
	db, exists := c.Get("db")
	if !exists {
		c.JSON(500, gin.H{"error": "Database not found"})
		return nil
	}
	return db.(*engine.Database)
}
