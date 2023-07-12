package controllers

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"time"

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

	err = bc.tpl.ExecuteTemplate(c.Writer, "buckets.gohtml", result)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}
}

type S3FileInfo struct {
	Key          string
	LastModified *time.Time
	Size         int64
}

func (bc *BucketsController) ListObjectsOfBucket(c *gin.Context) {
	bucketName := c.Query("name")
	ctx := context.TODO()
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	}

	result, err := bc.s3Client.ListObjectsV2(ctx, input)

	// if 301 is returned, most likely is that the bucket belongs to other region
	if err != nil {
		c.Error(err)
		c.Errors.JSON()
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, c.Errors.Errors())
		log.Print(err)
		return
	}

	fileInfo := S3FileInfo{
		Key:          aws.ToString(result.Contents[0].Key),
		LastModified: result.Contents[0].LastModified,
		Size:         result.Contents[0].Size,
	}

	c.JSON(http.StatusOK, fileInfo)

}
