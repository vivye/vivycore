package gS3

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"os"
)

func Upload(sess *session.Session, bucket string, filename string, data []byte) (err error) {
	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	} else {
		fmt.Printf("file uploaded to %s\n", result.Location)
		return nil
	}
}

func Download(sess *session.Session, bucket string, filename string, downloadFilename string) (err error) {
	downloader := s3manager.NewDownloader(sess)

	f, err := os.Create(downloadFilename)
	if err != nil {
		return fmt.Errorf("failed to create file %q, %v", filename, err)
	}

	n, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return fmt.Errorf("failed to download file, %v", err)
	} else {
		fmt.Printf("file downloaded, %d bytes\n", n)
		return nil
	}
}

func Read(sess *session.Session, bucket string, filename string) (ioReadCloser io.ReadCloser, err error) {
	client := NewS3Client(sess)
	result, err := client.GetObjectWithContext(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return nil, fmt.Errorf("Couldn't get object %v:%v. Here's why: %v\n", bucket, filename, err)
	}
	return result.Body, nil
}
