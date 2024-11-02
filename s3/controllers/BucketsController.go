package controllers

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (bc *BucketsController) ListBuckets(c *gin.Context) {
	ctx := context.TODO()
	input := &s3.ListBucketsInput{}

	result, err := bc.s3Client.ListBuckets(ctx, input)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})

		return
	}

	err = bc.tpl.ExecuteTemplate(c.Writer, "buckets.html", result)
	if err != nil {
		log.Print(err)
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

	result, err := bc.s3Client.ListObjectsV2(ctx, input) // todo which context should we pass?
	if err != nil {                                      // if 301 is returned, the bucket is probably in other region
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func (bc *BucketsController) UploadFile(c *gin.Context) {
	ctx := context.TODO()

	input := &s3.ListBucketsInput{}

	result, err := bc.s3Client.ListBuckets(ctx, input)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})

		return
	}

	err = bc.tpl.ExecuteTemplate(c.Writer, "upload.html", result)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}
}

func (bc *BucketsController) Upload(c *gin.Context) {
	bucketName := c.PostForm("bucket")
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}
	fileName := uuid.NewString()
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}
	ctx := context.TODO()

	input := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   file,
	}

	_, err = bc.s3Client.PutObject(ctx, input)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"resourceCreated": fmt.Sprintf("%s:%s", bucketName, fileName),
	})
}
