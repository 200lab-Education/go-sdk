package aws

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	s32 "github.com/aws/aws-sdk-go/service/s3"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func (s *s3) Upload(ctx context.Context, fileName string, cloudFolder string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size)

	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	ext := filepath.Ext(file.Name())
	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	fileKey := fmt.Sprintf("/%s/%s", cloudFolder, newFileName)
	params := &s32.PutObjectInput{
		Bucket:        aws.String(s.cfg.s3Bucket),
		Key:           aws.String(fileKey),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}

	_, err = s.service.PutObjectWithContext(ctx, params)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s.s3.amazonaws.com%s", s.cfg.s3Bucket, fileKey), nil
}

func (s *s3) UploadFileData(ctx context.Context, fileData []byte, fileName string) (string, error) {
	fileBytes := bytes.NewReader(fileData)
	fileType := http.DetectContentType(fileData)

	fileKey := fmt.Sprintf("/%s", fileName)
	params := &s32.PutObjectInput{
		Bucket:        aws.String(s.cfg.s3Bucket),
		Key:           aws.String(fileKey),
		Body:          fileBytes,
		ContentLength: aws.Int64(int64(len(fileData))),
		ContentType:   aws.String(fileType),
	}

	_, err := s.service.PutObjectWithContext(ctx, params)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s.s3.amazonaws.com%s", s.cfg.s3Bucket, fileKey), nil
}
