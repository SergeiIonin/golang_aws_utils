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
	"github.com/google/uuid"

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

	result, err := bc.s3Client.ListObjectsV2(ctx, input)

	// if 301 is returned, most likely is that the bucket belongs to other region
	if err != nil {
		c.Error(err)
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
		c.Error(err)
		c.JSON(http.StatusBadRequest, c.Errors.JSON())
		return
	}
	fileName := uuid.NewString()
	file, err := fileHeader.Open()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, c.Errors.JSON())
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
		c.Error(err)
		log.Print(err)
		c.JSON(http.StatusInternalServerError, c.Errors.JSON())
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"resourceCreated": fmt.Sprintf("%s:%s", bucketName, fileName),
	})
}
