package services

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	s3session *s3.S3
)

func init() {
	s3session = s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String("sa-east-1"),
	})))
}

func ListBuckets() (resp *s3.ListBucketsOutput) {
	resp, err := s3session.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		panic(err)
	}

	return resp
}

func CreateBucket(bucket_name string, region string) (resp *s3.CreateBucketOutput) {
	resp, err := s3session.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket_name),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(region),
		},
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				fmt.Println(s3.ErrCodeBucketAlreadyExists, aerr.Error())
				panic(err)
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				fmt.Println(s3.ErrCodeBucketAlreadyOwnedByYou, aerr.Error())
				panic(err)
			default:
				panic(err)
			}
		}
	}
	return resp
}

func UploadFile(filename string, bucket_name string) (resp *s3.PutObjectOutput) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	fmt.Println("Uploading file to S3 bucket:", filename)
	resp, err = s3session.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket_name),
		Key:    aws.String(filename),
		Body:   f,
		ACL:    aws.String(s3.BucketCannedACLPublicRead),
	})

	if err != nil {
		panic(err)
	}

	return resp

}

func ListObjects(bucket_name string) (resp *s3.ListObjectsV2Output) {
	resp, err := s3session.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket_name),
	})
	if err != nil {
		panic(err)
	}

	return resp
}

func GetFile(bucket_name string, filename string) {
	fmt.Println("Downloading file from S3 bucket:", filename)
	resp, err := s3session.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket_name),
		Key:    aws.String(filename),
	})
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	err = ioutil.WriteFile(filename, body, 0644)
	if err != nil {
		panic(err)
	}
}
