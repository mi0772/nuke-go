package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mi0772/nuke-go/engine"
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

func Pop(c *gin.Context) {
	database := getDatabase(c)
	key := c.Param("key")

	item, err := database.Pop(key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"item": item,
	})
}

func Read(c *gin.Context) {
	database := getDatabase(c)
	key := c.Param("key")

	item, err := database.Read(key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"item": item,
	})
}

func PushFile(c *gin.Context) {
	database := getDatabase(c)
	// Il metodo FormFile gestisce il file caricato
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}

	// Apri il file caricato
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer f.Close()

	// Leggi il contenuto del file in memoria
	fileBytes, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ora il file è memorizzato in fileBytes come []byte
	// Puoi fare ciò che desideri con fileBytes, ad esempio salvarlo in un database, elaborarlo, ecc.

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

	// In questo esempio, rispondiamo solo con la dimensione del file
	c.JSON(http.StatusOK, gin.H{
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
