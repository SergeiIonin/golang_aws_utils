package main

import (
	"GoAWSs3Utils/config"
	"GoAWSs3Utils/controllers"

	"html/template"

	"github.com/gin-gonic/gin"
)

var controller *controllers.BucketsController

func init() {
	// get s3 client
	s3Client := config.GetS3Client()
	tpl := template.Must(template.ParseGlob("templates/*.gohtml"))
	// create controller
	controller = controllers.NewBucketsController(s3Client, tpl)
}

func main() {
	router := gin.Default()
	router.GET("/buckets", controller.ListBuckets)
	// router.GET("/buckets/:name", getBucketByName)
	// router.POST("/buckets", postFileToBucket)

	router.Run("localhost:8080")
}
