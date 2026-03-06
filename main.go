package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"
	"github.com/gin-gonic/gin"
)

func main() {
	var queueClient *azqueue.QueueClient

	connStr := os.Getenv("AZURE_STORAGE_CONNECTION_STRING")
	queueName := os.Getenv("AZURE_STORAGE_QUEUE_NAME")
	if connStr != "" && queueName != "" {
		var err error
		queueClient, err = azqueue.NewQueueClientFromConnectionString(connStr, queueName, nil)
		if err != nil {
			log.Fatalf("failed to create queue client: %v", err)
		}
	} else {
		log.Println("AZURE_STORAGE_CONNECTION_STRING or AZURE_STORAGE_QUEUE_NAME not set; /message will return 503")
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/message", func(c *gin.Context) {
		if queueClient == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "queue not configured"})
			return
		}

		msg := struct {
			Message   string    `json:"message"`
			Timestamp time.Time `json:"timestamp"`
		}{
			Message:   "hello world",
			Timestamp: time.Now().UTC(),
		}

		b, err := json.Marshal(msg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal message"})
			return
		}

		if _, err = queueClient.EnqueueMessage(context.Background(), string(b), nil); err != nil {
			log.Printf("EnqueueMessage error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, msg)
	})

	r.Run(":8080")
}
