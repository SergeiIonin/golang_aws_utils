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

// ListBuckets lists all buckets in the s3 client and return an html page
func (bc *BucketsController) ListBuckets(c *gin.Context) {
	ctx := context.TODO()
	input := &s3.ListBucketsInput{}

	result, err := bc.s3Client.ListBuckets(ctx, input)
	if err != nil {
		log.Fatal(err)
		// todo switch on error code and return a status code
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	buckets := make([]Bucket, len(result.Buckets))
	for i, b := range result.Buckets {
		buckets[i] = Bucket{
			Name: aws.ToString(b.Name),
		}
	}

	err = bc.tpl.ExecuteTemplate(c.Writer, "buckets.gohtml", buckets)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}
}
