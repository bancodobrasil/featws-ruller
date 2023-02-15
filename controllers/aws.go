package controllers

import (
	"net/http"

	"github.com/bancodobrasil/featws-ruller/services"
	"github.com/gin-gonic/gin"
)

func ListBucketsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get a list of buckets from S3
		resp := services.ListBuckets()

		// Return the buckets as a JSON response
		c.JSON(http.StatusOK, gin.H{"buckets": resp.Buckets})
	}
}

func CreateBucketHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the bucket name and region from the request body
		var reqBody struct {
			BucketName string `json:"bucket_name" binding:"required"`
			Region     string `json:"region" binding:"required"`
		}

		if err := c.BindJSON(&reqBody); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Create a new S3 bucket
		resp := services.CreateBucket(reqBody.BucketName, reqBody.Region)

		// Return the response as a JSON response
		c.JSON(http.StatusOK, gin.H{"response": resp})
	}
}

func UploadFileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the file name and bucket name from the request body
		var reqBody struct {
			FileName   string `json:"file_name" binding:"required"`
			BucketName string `json:"bucket_name" binding:"required"`
		}

		if err := c.BindJSON(&reqBody); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Upload the file to the specified S3 bucket
		resp := services.UploadFile(reqBody.FileName, reqBody.BucketName)

		// Return the response as a JSON response
		c.JSON(http.StatusOK, gin.H{"response": resp})
	}
}

func ListObjectsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the bucket name from the request URL
		bucketName := c.Param("bucket_name")

		// List the objects in the specified S3 bucket
		resp := services.ListObjects(bucketName)

		// Return the objects as a JSON response
		c.JSON(http.StatusOK, gin.H{"objects": resp.Contents})
	}
}

func GetFileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the bucket name and file name from the request URL
		bucketName := c.Param("bucket_name")
		fileName := c.Param("file_name")

		// Download the file from the specified S3 bucket
		services.GetFile(bucketName, fileName)

		// Return a success message as a JSON response
		c.JSON(http.StatusOK, gin.H{"message": "File downloaded successfully"})
	}
}
