package controllers

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

type BucketsController struct {
	s3Client *s3.Client
	tpl      *template.Template
}

func NewBucketsController(client *s3.Client, tpl *template.Template) *BucketsController {
	return &BucketsController{
		s3Client: client,
		tpl:      tpl,
	}
}

type Bucket struct {
	Name string
}

func (controller BucketsController) ListBuckets(ctx *gin.Context) {
	bucketsOutput, err := controller.s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatal(err)
		// todo switch on error code and return a status code
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	buckets := make([]Bucket, 0, 10)

	for _, bucket := range bucketsOutput.Buckets {
		buckets = append(buckets, Bucket{Name: aws.ToString(bucket.Name)})
	}

	// create file

	// file, _ := os.Create("templates/bucketsRes.html")
	// controller.tpl.ExecuteTemplate(file, "buckets.gohtml", buckets)

	// ctx.JSON(http.StatusOK, gin.H{
	// 	"buckets": buckets,
	// })

	// return html template
	ctx.HTML(http.StatusOK, "buckets.gohtml", buckets)

}
