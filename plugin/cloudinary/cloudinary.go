package cloudinary

import (
	"context"
	"flag"
	"fmt"
	"github.com/200Lab-Education/go-sdk/logger"
	"github.com/200Lab-Education/go-sdk/sdkcm"
)

var (
	ErrCloudinaryApiKeyMissing       = sdkcm.CustomError("ErrCloudinaryApiKeyMissing", "Cloudinary API key is missing")
	ErrCloudinaryApiSecretKeyMissing = sdkcm.CustomError("ErrCloudinaryApiSecretKeyMissing", "Cloudinary API secret key is missing")
	ErrCloudinaryCloudNameMissing    = sdkcm.CustomError("ErrCloudinaryApiKeyMissing", "Cloudinary cloud name is missing")
)

type Cloudinary interface {
	VideoUpload(ctx context.Context, file string, uploadPreset string, folder string, format string) (*VideoResult, error)
}

type cloudinary struct {
	name   string
	prefix string
	logger logger.Logger

	config cloudinaryConfig
}

type cloudinaryConfig struct {
	apiKey    string
	apiSecret string
	cloudName string
}

func New(prefix ...string) *cloudinary {
	pre := "cloudinary"

	if len(prefix) > 0 {
		pre = prefix[0]
	}

	return &cloudinary{
		name:   "cloudinary",
		prefix: pre,
	}
}

func (cd *cloudinary) Get() interface{} {
	return cd
}

func (cd *cloudinary) Name() string {
	return cd.name
}

func (cd *cloudinary) InitFlags() {
	flag.StringVar(&cd.config.apiKey, fmt.Sprintf("%s-%s", cd.GetPrefix(), "api-key"), "", "Cloudinary api key")
	flag.StringVar(&cd.config.apiSecret, fmt.Sprintf("%s-%s", cd.GetPrefix(), "api-secret"), "", "Cloudinary api secret")
	flag.StringVar(&cd.config.cloudName, fmt.Sprintf("%s-%s", cd.GetPrefix(), "cloud-name"), "", "Cloudinary cloud name")
}

func (cd *cloudinary) Configure() error {
	cd.logger = logger.GetCurrent().GetLogger(cd.Name())

	if err := cd.config.check(); err != nil {
		cd.logger.Errorln(err)
		return err
	}

	return nil
}

func (cd *cloudinary) Run() error {
	return cd.Configure()
}

func (cd *cloudinary) Stop() <-chan bool {
	c := make(chan bool)
	go func() { c <- true }()
	return c
}

func (cd *cloudinary) GetPrefix() string {
	return cd.prefix
}

func (cfg *cloudinaryConfig) check() error {
	if len(cfg.apiKey) < 1 {
		return ErrCloudinaryApiKeyMissing
	}
	if len(cfg.apiSecret) < 1 {
		return ErrCloudinaryApiSecretKeyMissing
	}
	if len(cfg.cloudName) < 1 {
		return ErrCloudinaryCloudNameMissing
	}
	return nil
}
