package main

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const BUCKET_NAME = "therion"

type S3Sync struct {
	client *s3.Client
	bucket *string
}

func NewS3Sync() *S3Sync {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("therion-user"))
	if err != nil {
		panic(err)
	}
	var bucket string = BUCKET_NAME
	return &S3Sync{client: s3.NewFromConfig(cfg), bucket: &bucket}
}

func (s S3Sync) getFiles() []file {
	var files []file
	response, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{Bucket: s.bucket})
	if err != nil {
		panic(err)
	}
	for _, object := range response.Contents {
		files = append(files, file{path: *object.Key, modTime: *object.LastModified, data: []byte{}})
	}
	return files
}

func (s S3Sync) uploadFile(f *file) {
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{Bucket: s.bucket, Key: &f.path, Body: bytes.NewReader(f.data)})
	if err != nil {
		panic(err)
	}
}

func (s S3Sync) downloadFile(f *file) {
}
