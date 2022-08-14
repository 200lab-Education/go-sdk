/*----------------------------------------------------------------*\
 * @author          Ly Nam <lyquocnam@live.com>
 * @copyright       2019 Viet Tran <viettranx@gmail.com>
 * @license         Apache-2.0
 * @description		Plugin to upload image to AWS S3
 *----------------------------------------------------------------*/
package aws

import (
	"flag"
	"fmt"
	"github.com/200Lab-Education/go-sdk/logger"
	"github.com/200Lab-Education/go-sdk/sdkcm"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	s32 "github.com/aws/aws-sdk-go/service/s3"
)

var (
	ErrS3ApiKeyMissing       = sdkcm.CustomError("ErrS3ApiKeyMissing", "AWS S3 API key is missing")
	ErrS3ApiSecretKeyMissing = sdkcm.CustomError("ErrS3ApiSecretKeyMissing", "AWS S3 API secret key is missing")
	ErrS3RegionMissing       = sdkcm.CustomError("ErrS3RegionMissing", "AWS S3 region is missing")
	ErrS3BucketMissing       = sdkcm.CustomError("ErrS3ApiKeyMissing", "AWS S3 bucket is missing")
)

type s3 struct {
	name   string
	prefix string
	logger logger.Logger

	cfg s3Config

	session *session.Session
	service *s32.S3
}

type s3Config struct {
	s3ApiKey    string
	s3ApiSecret string
	s3Region    string
	s3Bucket    string
}

func New(prefix ...string) *s3 {
	pre := "aws-s3"

	if len(prefix) > 0 {
		pre = prefix[0]
	}

	return &s3{
		name:   "aws-s3",
		prefix: pre,
	}
}

func (s *s3) Get() interface{} {
	return s
}

func (s *s3) Name() string {
	return s.name
}

func (s *s3) InitFlags() {
	flag.StringVar(&s.cfg.s3ApiKey, fmt.Sprintf("%s-%s", s.GetPrefix(), "api-key"), "", "S3 API key")
	flag.StringVar(&s.cfg.s3ApiSecret, fmt.Sprintf("%s-%s", s.GetPrefix(), "api-secret"), "", "S3 API secret key")
	flag.StringVar(&s.cfg.s3Region, fmt.Sprintf("%s-%s", s.GetPrefix(), "region"), "", "S3 region")
	flag.StringVar(&s.cfg.s3Bucket, fmt.Sprintf("%s-%s", s.GetPrefix(), "bucket"), "", "S3 bucket")
}

func (s *s3) Configure() error {
	s.logger = logger.GetCurrent().GetLogger(s.Name())

	if err := s.cfg.check(); err != nil {
		s.logger.Errorln(err)
		return err
	}

	credential := credentials.NewStaticCredentials(s.cfg.s3ApiKey, s.cfg.s3ApiSecret, "")
	_, err := credential.Get()
	if err != nil {
		s.logger.Errorln(err)
		return err
	}

	config := aws.NewConfig().WithRegion(s.cfg.s3Region).WithCredentials(credential)
	ss, err := session.NewSession(config)
	service := s32.New(ss, config)

	s.session = ss
	s.service = service

	return nil
}

func (s *s3) GetPrefix() string {
	return s.prefix
}

func (s *s3) Run() error {
	return s.Configure()
}

func (s *s3) Stop() <-chan bool {
	c := make(chan bool)
	go func() { c <- true }()
	return c
}

func (cfg *s3Config) check() error {
	if len(cfg.s3ApiKey) < 1 {
		return ErrS3ApiKeyMissing
	}
	if len(cfg.s3ApiSecret) < 1 {
		return ErrS3ApiSecretKeyMissing
	}
	if len(cfg.s3Bucket) < 1 {
		return ErrS3BucketMissing
	}
	if len(cfg.s3Region) < 1 {
		return ErrS3RegionMissing
	}
	return nil
}
