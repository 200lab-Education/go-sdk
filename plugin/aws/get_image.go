package aws

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	s32 "github.com/aws/aws-sdk-go/service/s3"
	"time"
)

func (s *s3) GetImageWithExpireLink(ctx context.Context, imageKey string, duration time.Duration) (string, error) {
	req, _ := s.service.GetObjectRequest(&s32.GetObjectInput{
		Bucket: aws.String(s.cfg.s3Bucket),
		Key:    aws.String(imageKey),
	})

	return req.Presign(duration)
}
