package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
)

// LoggerMiddleware membuat log lengkap dengan Request & Response format custom
func LoggerMiddleware() gin.HandlerFunc {
	// Buka file log
	logFile, err := os.OpenFile("backend.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Gagal membuka file log: %v", err)
	}
	logger := log.New(logFile, "", 0)

	return func(c *gin.Context) {
		// Generate unique request ID
		requestID := uuid.New().String()
		startTime := time.Now()

		// Baca body request
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // supaya bisa dibaca ulang oleh Gin

		// Log START
		logger.Printf("[Start][RequestId]= %s, [Path]= %s, [IP]= %s, [Time]= %s, [Request]= %s\n",
			requestID,
			c.FullPath(),
			c.ClientIP(),
			startTime.Format("2006-01-02 15:04:05.000000"),
			string(bodyBytes),
		)

		// Tangkap response writer
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next() // Jalankan handler

		stopTime := time.Now()
		responseBody := blw.body.String()

		// Format JSON response agar rapi
		var formattedJSON bytes.Buffer
		if json.Valid([]byte(responseBody)) {
			json.Indent(&formattedJSON, []byte(responseBody), "", "  ")
		} else {
			formattedJSON.WriteString(responseBody)
		}

		// Log STOP
		logger.Printf("[Stop][RequestId]= %s, [Path]= %s, [IP]= %s, [Time]= %s, [Response]= %s\n\n",
			requestID,
			c.FullPath(),
			c.ClientIP(),
			stopTime.Format("2006-01-02 15:04:05.000000"),
			formattedJSON.String(),
		)
	}
}

// bodyLogWriter digunakan untuk menangkap isi response Gin
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
