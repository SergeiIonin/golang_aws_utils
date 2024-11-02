package main

import (
	"html/template"
	"log"

	"GoAWSs3Utils/s3/config"
	"GoAWSs3Utils/s3/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	s3Client := config.GetS3Client()
	tpl := template.Must(template.ParseGlob("templates/*.html"))
	controller := controllers.NewBucketsController(s3Client, tpl)

	router := gin.Default()
	router.GET("/buckets", controller.ListBuckets)
	router.GET("/buckets/objects", controller.ListObjectsOfBucket)
	router.GET("/buckets/uploadFiles", controller.UploadFile)
	router.POST("/buckets/upload", controller.Upload)

	if err := router.Run("localhost:8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
