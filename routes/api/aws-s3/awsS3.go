package awss3

import (
	"github.com/bancodobrasil/featws-ruller/controllers"
	"github.com/gin-gonic/gin"
)

func awsS3Router(router *gin.RouterGroup) {
	router.POST("/create-bucket-s3", controllers.CreateBucketHandler())
	router.POST("/upload-file-s3", controllers.UploadFileHandler())
	router.GET("/list-buckets-s3", controllers.ListBucketsHandler())
	router.GET("/list-objects-s3/:bucket_name", controllers.ListObjectsHandler())
	router.GET("/download-file-s3/:bucket_name/:file_name", controllers.GetFileHandler())

}
