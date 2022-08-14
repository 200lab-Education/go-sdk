package aws

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	s32 "github.com/aws/aws-sdk-go/service/s3"
)

func (s *s3) DeleteImages(ctx context.Context, fileKeys []string) error {

	del := &s32.Delete{
		Objects: toOIDs(fileKeys),
		Quiet:   aws.Bool(false),
	}

	doi := &s32.DeleteObjectsInput{
		Bucket: aws.String(s.cfg.s3Bucket),
		Delete: del,
	}

	res, err := s.service.DeleteObjects(doi)
	if err != nil {
		return err
	}

	s.logger.Infoln(res)

	return nil
}

func toOIDs(keys []string) []*s32.ObjectIdentifier {
	ret := make([]*s32.ObjectIdentifier, len(keys))
	for i := 0; i < len(ret); i++ {
		oid := &s32.ObjectIdentifier{
			Key: &(keys[i]),
		}
		ret[i] = oid
	}
	return ret
}

func (s *s3) DeleteObject(ctx context.Context, key string) error {
	input := &s32.DeleteObjectInput{
		Bucket: aws.String(s.cfg.s3Bucket),
		Key:    aws.String(key),
	}

	res, err := s.service.DeleteObject(input)
	if err != nil {
		return err
	}

	s.logger.Infoln(res)

	return nil
}
