package main

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type file struct {
	name    string
	modTime time.Time
}

type Syncer interface {
	getFiles() []file
	uploadFiles([]file)
}

type S3Sync struct {
	client *s3.Client
}

func NewS3Sync() *S3Sync {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("therion-user"))
	if err != nil {
		panic(err)
	}
	return &S3Sync{client: s3.NewFromConfig(cfg)}
}

func (s S3Sync) getFiles() []file {
	return nil
}

func (s S3Sync) uploadFiles(files []file) {

}

func sync(path string, syncer Syncer) {
	var localFiles []file
	err := filepath.WalkDir(path, func(name string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".md" {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			localFiles = append(localFiles, file{name: entry.Name(), modTime: info.ModTime()})
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
