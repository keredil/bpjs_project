package main

// @title BPJS API
// @version 1.0
// @description API untuk menampilkan program BPJS Ketenagakerjaan
// @host localhost:8080
// @BasePath /

import (
	"backend/config"
	"backend/routes"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "backend/docs" // generated docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// ==== Generate UUID-like Request ID ====
func generateRequestID() string {
	const charset = "abcdef0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 36)
	for i := range b {
		switch i {
		case 8, 13, 18, 23:
			b[i] = '-'
		default:
			b[i] = charset[rand.Intn(len(charset))]
		}
	}
	return string(b)
}

// ==== Body capture untuk response ====
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// ==== Custom Logger mirip API BPJS ====
func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		requestID := generateRequestID()
		path := c.Request.URL.Path
		ip := c.ClientIP()

		// Baca body request (untuk POST/PUT)
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		log.Printf("[Start][RequestId]= %s, [Path]= %s, [IP]= %s, [Time]= %s, [Request]= %s\n",
			requestID,
			path,
			ip,
			startTime.Format("2006-01-02 15:04:05.000000"),
			string(requestBody),
		)

		// Tangkap response
		responseBody := &bytes.Buffer{}
		writer := &bodyLogWriter{body: responseBody, ResponseWriter: c.Writer}
		c.Writer = writer

		c.Next()

		stopTime := time.Now()

		// Format response JSON agar rapi
		formattedResponse := responseBody.String()
		var prettyJSON bytes.Buffer
		if json.Valid(responseBody.Bytes()) {
			_ = json.Indent(&prettyJSON, responseBody.Bytes(), "", "  ")
			formattedResponse = prettyJSON.String()
		}

		log.Printf("[Stop][RequestId]= %s, [Path]= %s, [IP]= %s, [Time]= %s, [Response]= %s\n\n",
			requestID,
			path,
			ip,
			stopTime.Format("2006-01-02 15:04:05.000000"),
			formattedResponse,
		)
	}
}

func main() {
	// === Setup logging ke file backend.log ===
	logFile, err := os.OpenFile("backend.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Tidak dapat membuka atau membuat file log:", err)
	}

	// === Tulis log ke file dan terminal (stdout) ===
	multiWriter := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(multiWriter)

	r := gin.New()        // gunakan gin.New() agar tidak double middleware
	r.Use(gin.Recovery()) // middleware default recovery
	r.Use(CustomLogger()) // custom logger

	// === Middleware CORS ===
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// === Koneksi database dan routes ===
	config.ConnectDatabase()
	routes.SetupRoutes(r)

	// === Swagger ===
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// === Jalankan server ===
	fmt.Println("🚀 Server berjalan di http://localhost:8080")
	r.Run(":8080")
}
