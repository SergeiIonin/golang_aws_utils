package main

import (
	"GoAWSs3Utils/config"
	"GoAWSs3Utils/controllers"

	"html/template"

	"github.com/gin-gonic/gin"
)

var controller *controllers.BucketsController

func init() {
	s3Client := config.GetS3Client()
	tpl := template.Must(template.ParseGlob("templates/*.html"))
	controller = controllers.NewBucketsController(s3Client, tpl)
}

func main() {
	router := gin.Default()
	router.GET("/buckets", controller.ListBuckets)
	router.GET("/buckets/objects", controller.ListObjectsOfBucket)
	router.GET("/buckets/uploadFiles", controller.UploadFile)
	router.POST("/buckets/upload", controller.Upload)

	router.Run("localhost:8080")
}
